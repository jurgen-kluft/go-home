package main

import (
	"net"
	"sync"
)

// Usually the MTU is <1500 bytes
const mtu = 1500

// Connect dials to address and returns a send-only
// channel and an error.
func Connect(address string) (chan<- []byte, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	outbound := make(chan []byte, 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Caught a panic, so we're going to close
				// the outbound channel. Receivers are supposed
				// to check for closed channels, so we're good.
				close(outbound)
			}
		}()
		for message := range outbound {
			conn.Write(message)
		}
	}()

	return outbound, nil
}

// Listen starts a UDP listener at address and
// returns a read-only channel and an error.
// If a value is sent to close, the UDP socket will
// be closed.
func Listen(address string, close chan bool) (<-chan []byte, *PacketPool, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, nil, err
	}

	inbound := make(chan []byte, 10)
	inpackets := NewPacketPool()
	go func() {
		for {
			select {
			case <-close:
				conn.Close()
				return
			default:
			}

			b := inpackets.Get()
			n, err := conn.Read(b)
			if err != nil {
				continue
			}

			b = b[:n]
			inbound <- b
		}
	}()

	return inbound, inpackets, nil
}

// PacketPool is a sync.Pool for bytes.Buffer objects
type PacketPool struct {
	pool sync.Pool
}

// NewPacketPool creates a new PacketPool instance
func NewPacketPool() *PacketPool {
	var pp PacketPool
	pp.pool.New = allocPacket
	return &pp
}

func allocPacket() interface{} {
	return make([]byte, mtu)
}

// Get returns a bytes.Buffer from the specified pool
func (pp *PacketPool) Get() []byte {
	return pp.pool.Get().([]byte)
}

// Release puts the given bytes.Buffer back in the specified pool after
// resetting it
func (pp *PacketPool) Release(buf []byte) {
	pp.pool.Put(buf)
}
