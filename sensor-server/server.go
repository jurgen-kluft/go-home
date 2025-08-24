package main

// TODO:
// This is a very simple TCP server that listens for messages from ESP32 devices.
// It receives messages in a custom binary format containing different types of sensor data.
// The server decodes the messages and generates MQTT messages to be sent to an MQTT broker.
// We also will push this data to an InfluxDB database that we can visualize with Grafana.
// The server will be able to handle multiple clients and will be able to process messages in parallel.

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jurgen-kluft/go-home/sensor-server/hollywood/actor"
)

type handler struct{}

func newHandler() actor.Receiver {
	return &handler{}
}

func (handler) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case []byte:
		fmt.Println("received packet")

		// Custom binary message format
		packet, err := DecodeNetworkPacket(msg)
		if err != nil {
			slog.Error("decode error", "err", err)
			return
		}
		fmt.Println("decoded packet:", packet)

	case actor.Stopped:
		for i := 0; i < 3; i++ {
			fmt.Printf("\r handler stopping in %d", 3-i)
			time.Sleep(time.Second)
		}
		fmt.Println("handler stopped")
	}
}

type session struct {
	conn net.Conn
}

func newSession(conn net.Conn) actor.Producer {
	return func() actor.Receiver {
		return &session{
			conn: conn,
		}
	}
}

func (s *session) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Initialized:
	case actor.Started:
		slog.Info("new connection", "addr", s.conn.RemoteAddr())
		go s.readLoop(c)
	case actor.Stopped:
		s.conn.Close()
	}
}

func (s *session) readLoop(c *actor.Context) {
	for {
		// What is the maximum size of a message?
		msg := make([]byte, 1024)
		n, err := s.conn.Read(msg)
		if err != nil {
			slog.Error("conn read error", "err", err)
			break
		}

		// Send to the handler to process the message
		c.Send(c.Parent().Child("handler/default"), msg[:n])
	}
	// Loop is done due to error or we need to close due to server shutdown.
	c.Send(c.Parent(), &connRem{pid: c.PID()})
}

type connAdd struct {
	pid  *actor.PID
	conn net.Conn
}

type connRem struct {
	pid *actor.PID
}

type server struct {
	listenAddr string
	ln         net.Listener
	sessions   map[*actor.PID]net.Conn
}

func newServer(listenAddr string) actor.Producer {
	return func() actor.Receiver {
		return &server{
			listenAddr: listenAddr,
			sessions:   make(map[*actor.PID]net.Conn),
		}
	}
}

func (s *server) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		ln, err := net.Listen("tcp", s.listenAddr)
		if err != nil {
			panic(err)
		}
		s.ln = ln
		// start the handler that will handle the incomming messages from clients/sessions.
		c.SpawnChild(newHandler, "handler", actor.WithID("default"))
	case actor.Started:
		slog.Info("server started", "addr", s.listenAddr)
		go s.acceptLoop(c)
	case actor.Stopped:
		// on stop all the childs sessions will automatically get the stop
		// message and close all their underlying connection.
	case *connAdd:
		slog.Debug("added new connection to my map", "addr", msg.conn.RemoteAddr(), "pid", msg.pid)
		s.sessions[msg.pid] = msg.conn
	case *connRem:
		slog.Debug("removed connection from my map", "pid", msg.pid)
		delete(s.sessions, msg.pid)
	}
}

func (s *server) acceptLoop(c *actor.Context) {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			break
		}
		pid := c.SpawnChild(newSession(conn), "session", actor.WithID(conn.RemoteAddr().String()))
		c.Send(c.PID(), &connAdd{
			pid:  pid,
			conn: conn,
		})
	}
}

func main() {
	listenAddr := flag.String("listen", "10.0.0.66:31337", "listen address of the TCP server")

	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}
	serverPID := e.Spawn(newServer(*listenAddr), "server")

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	<-sigch

	<-e.Poison(serverPID).Done()
}
