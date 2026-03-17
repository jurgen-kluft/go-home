# micro-service

Lightweight local micro-service messaging over Unix domain sockets.

This package implements a small framed message protocol and connection management for local processes to communicate via Unix domain sockets. It is designed for reliability and low overhead: automatic reconnect with backoff, message pooling, broadcasting, and simple readiness checks.

**Highlights**
- Framed binary protocol with a 32-byte header (magic 'UMQ1', version 1).
- Max payload 64 KiB and pooled buffers for low GC pressure.
- Automatic dialer/acceptor model: services listen on a socket path and can dial peers.
- Convenience helpers: `NewMessage`, `NewTextMessage`, `SendTo`, `Broadcast`.
- Readiness and peer management: `Ready`, `ReadyWait`, `Peers`, `Connect`, `Disconnect`.
- Hooks: `Config.OnConnect` and `Config.OnDisconnect` for observability.

Quick overview

- Package: `microservice`
- Constructor: `New(name string, peerPath string, tickFrequency time.Duration) (*Service, error)` — creates a listener at `peerPath` and starts accept/dial loops.
- Message type: `Message` with helpers `MsgType()`, `Len()`, `SrcID()`, `DstID()` and `Release()` for pooled frames.
- Delegate type: `type Delegate func(m *Service, message *Message) bool` — register handlers keyed by peer path or id using `RegisterHandler` (see source).

Configuration (partial)

- `Config.ListenPath` — filesystem path for the Unix socket the service listens on (required).
- `Config.Peers` — slice of peer socket paths to dial on startup.
- `Config.InboundCapacity`, `Config.OutboundCapacity` — channel capacities.
- `Config.ReconnectMin`, `Config.ReconnectMax` — backoff range for dial retries.
- `Config.DropWhenFull` — if true, outbound messages may be dropped when queues are full.
- `Config.OnConnect`, `Config.OnDisconnect` — optional callbacks invoked with `(peerPath string, peerID uint32)`.

Usage example

```go
srv, err := microservice.New("myservice", "/tmp/myservice.sock", 5*time.Second)
if err != nil {
	// handle error
}

// register a handler for incoming messages from a specific peer path
srv.RegisterHandler("/tmp/other.sock", func(m *microservice.Service, msg *microservice.Message) bool {
	// inspect msg.MsgType(), msg.Payload etc.
	// return true to keep loop running, false to stop
	return true
})

// read inbound messages from the channel (or use registered handlers)
go func() {
	for msg := range srv.Inbound() {
		// process, then release pooled message
		// ...
		msg.Release()
	}
}()

// send a text message
_ = srv.SendJsonTo("/tmp/other.sock", `{"hello":"world"}`)

// shutdown
_ = srv.Close()
```

Notes

- The constructor removes any stale socket at the `ListenPath` and sets file mode to `0600`.
- Peer IDs are computed with `Hash32(path)`; these are used in message headers.
- The package exposes a tick message mechanism; the constructor seeds a tick into the inbound channel after `tickFrequency`.
- Buffer pooling is used; callers should call `Release()` on messages received from the inbound channel when done.

References

- Source: [micro-service/micro-service.go](micro-service/micro-service.go)
- Unix-specific implementation: [micro-service/micro-service-unix.go](micro-service/micro-service-unix.go)

Testing

Run package tests:

```
go test ./micro-service
```

If you want, I can add more examples or expand the API section with exact method signatures.

