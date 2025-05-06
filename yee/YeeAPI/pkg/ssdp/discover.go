package ssdp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/thomasf/lg"
	"github.com/thomasf/yeelight/pkg/yeel"
	"golang.org/x/net/ipv4"
)

func GetDevices(ifname string) ([]yeel.Device, error) {
	return GetDevicesWithContext(context.Background(), ifname)
}

func GetDevicesWithContext(ctx context.Context, ifname string) ([]yeel.Device, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	timeout := time.Second
	var devices []yeel.Device

	deviceCh := make(chan yeel.Device, 100)
	conn := Conn{timeout: timeout}
	if err := conn.Start(ifname); err != nil {
		return devices, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := conn.Listen(deviceCh)
		if err != nil {
			lg.Errorln(err)
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := conn.Search()
		if err != nil {
			lg.Errorln(err)
		}
	}()

	deviceMap := make(map[string]yeel.Device)

	timeoutT := time.NewTicker(timeout)

loop:
	for {
		select {
		case v := <-deviceCh:
			deviceMap[v.ID] = v
		case <-timeoutT.C:
			break loop
		case <-ctx.Done():
			break loop
		}
	}

	for _, v := range deviceMap {
		devices = append(devices, v)
	}

	time.Sleep(2 * time.Second)
	if err := conn.Stop(); err != nil {
		lg.Fatal(err)
	}

	wg.Wait()

	return devices, nil
}

const (
	ssdpDiscover = `"ssdp:discover"`
	// ntsAlive       = `ssdp:alive`
	// ntsByebye      = `ssdp:byebye`
	// ntsUpdate      = `ssdp:update`
	ssdpUDP4Addr = "239.255.255.250:1982"
	// ssdpSearchPort = 1982
	methodSearch = "M-SEARCH"
	// methodNotify   = "NOTIFY"
)

// Conn .
type Conn struct {
	uConn   net.PacketConn
	mConn   *ipv4.PacketConn
	timeout time.Duration
}

func (c *Conn) Stop() error {
	var errs []error
	if err := c.uConn.Close(); err != nil {
		lg.Warning(err)
		errs = append(errs, err)

	}
	if err := c.mConn.Close(); err != nil {
		lg.Warning(err)
		errs = append(errs, err)
	}

	return nil

}

func (c *Conn) Start(interfaceName string) error {
	netif, err := net.InterfaceByName(interfaceName)
	if err != nil {
		lg.Fatal(err)
	}

	group := net.IPv4(239, 255, 255, 250)

	uniConn, err := net.ListenPacket("udp4", "0.0.0.0:1982")
	if err != nil {
		return err
	}
	// defer c.Close()
	multiConn := ipv4.NewPacketConn(uniConn)

	// if err := p.JoinGroup(en0, &net.UDPAddr{IP: group, Port: 1982}); err != nil {
	if err := multiConn.JoinGroup(netif, &net.UDPAddr{IP: group}); err != nil {
		return err
	}

	if err := multiConn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		return err
	}

	c.uConn = uniConn
	c.mConn = multiConn
	return nil
}

func (c *Conn) Listen(deviceC chan yeel.Device) error {
loop:
	for {
		b := make([]byte, 2048)
		n, cm, src, err := c.mConn.ReadFrom(b)
		if err != nil {
			return err
		}
		if lg.V(10) {
			lg.Infoln(n, src, cm.Src, cm.Dst)
		}

		// pad the response with an extra \r\n to allow textproto ReadMIMEHeader() to be used.
		b[n] = '\r'
		b[n+1] = '\n'
		n = n + 2

		// Parse response.
		response, err := readResponse(bufio.NewReader(bytes.NewBuffer(b[:n])))
		if err != nil {
			if lg.V(5) {
				lg.Warningln("httpu: error while parsing response: %v", err)
			}
			if lg.V(10) {
				lg.Infoln(spew.Sdump(string(b[:n])))

			}
			continue loop
		}

		device, err := yeel.ParseDeviceFromHeader(response.Header)
		if err != nil {
			lg.Errorln(err)
			continue loop
		}

		deviceC <- device

	}
}

var req = http.Request{
	Method: methodSearch,
	Host:   ssdpUDP4Addr,
	URL:    &url.URL{Opaque: "*"},
	Header: http.Header{
		// Putting headers in here avoids them being title-cased.
		// (The UPnP discovery protocol uses case-sensitive headers)
		"HOST": []string{ssdpUDP4Addr},
		"MAN":  []string{ssdpDiscover},
		"ST":   []string{"wifi_bulb"},
	},
}

func (c *Conn) Search() error {
	timeout := 2 * time.Minute
	if c.timeout != 0 {
		timeout = c.timeout
	}
	err := c.do(&req, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) do(req *http.Request, timeout time.Duration) error {
	var requestBuf bytes.Buffer
	method := req.Method
	if method == "" {
		method = "GET"
	}
	if _, err := fmt.Fprintf(&requestBuf, "%s %s HTTP/1.1\r\n", method, req.URL.RequestURI()); err != nil {
		return err
	}
	if err := req.Header.Write(&requestBuf); err != nil {
		return err
	}
	if _, err := requestBuf.Write([]byte{'\r', '\n'}); err != nil {
		return err
	}
	destAddr, err := net.ResolveUDPAddr("udp", req.Host)
	if err != nil {
		return err
	}

	if n, err := c.uConn.WriteTo(requestBuf.Bytes(), destAddr); err != nil {
		return err
	} else if n < len(requestBuf.Bytes()) {
		return fmt.Errorf("httpu: wrote %d bytes rather than full %d in request",
			n, len(requestBuf.Bytes()))
	}
	return err
}

func readResponse(r *bufio.Reader) (*http.Response, error) {
	tp := textproto.NewReader(r)
	resp := &http.Response{}

	// Parse the first line of the response.
	line, err := tp.ReadLine()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	f := strings.SplitN(line, " ", 3)
	if len(f) < 2 {
		return nil, &badStringError{"malformed HTTP response", line}
	}
	reasonPhrase := ""
	if f[1] == "*" {
		if f[0] != "NOTIFY" {
			return nil, &badStringError{"malformed HTTP response", line}
		}
		resp.Status = "200"
	} else {
		if len(f) > 2 {
			reasonPhrase = f[2]
		}
		if len(f[1]) != 3 {
			return nil, &badStringError{"malformed HTTP status code", f[1]}
		}
		resp.StatusCode, err = strconv.Atoi(f[1])
		if err != nil || resp.StatusCode < 0 {
			return nil, &badStringError{"malformed HTTP status code", f[1]}
		}
		resp.Status = f[1] + " " + reasonPhrase
	}
	if f[1] == "*" {
		resp.Proto = f[2]
	} else {
		resp.Proto = f[0]
	}
	var ok bool
	if resp.ProtoMajor, resp.ProtoMinor, ok = http.ParseHTTPVersion(resp.Proto); !ok {
		return nil, &badStringError{"malformed HTTP version", resp.Proto}
	}

	// Parse the response headers.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	resp.Header = http.Header(mimeHeader)

	return resp, nil
}

type badStringError struct {
	what string
	str  string
}

func (e *badStringError) Error() string { return fmt.Sprintf("%s %q", e.what, e.str) }
