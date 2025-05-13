package main

// This is a very simple UDP server that listens for messages from ESP32 devices.
// It receives messages in a custom binary format containing different types of sensor data.
// The server decodes the messages and generates MQTT messages to be sent to an MQTT broker.
// We also will push this data to an InfluxDB database that we can visualize with Grafana.
// The server is designed to be efficient and will have to very stable.
// The server will be able to handle multiple clients and will be able to process messages in parallel.

// We will use `gnet` for the TCP/UDP framework

import (
	"flag"
	"fmt"
	"log"

	"github.com/jurgen-kluft/go-home/sensor-server/gnet"
)

type echoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

func main() {
	var port int
	var multicore bool

	// Example command: go run echo.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()
	echo := &echoServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: multicore}
	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))
}
