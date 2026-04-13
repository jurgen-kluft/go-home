//go:build darwin
// +build darwin

package microservice

import (
	"net"
	"syscall"
	"time"
)

func (s *Service) startDialer(peerPath string) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		peerID := Hash32(peerPath)
		backoff := s.cfg.ReconnectMin
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
			}
			d := net.Dialer{Timeout: 2 * time.Second, Control: setUnixNoSigPipe, KeepAlive: -1}
			conn, err := d.DialContext(s.ctx, "unix", peerPath)
			if err != nil {
				time.Sleep(jitter(backoff))
				if backoff < s.cfg.ReconnectMax {
					backoff *= 2
					if backoff > s.cfg.ReconnectMax {
						backoff = s.cfg.ReconnectMax
					}
				}
				continue
			}
			_ = conn.(*net.UnixConn).SetReadBuffer(MaxFrame * 2)
			_ = conn.(*net.UnixConn).SetWriteBuffer(MaxFrame * 2)
			h := &connHandle{parent: s, peerPath: peerPath, peerID: peerID, conn: conn.(*net.UnixConn), writeCh: make(chan *Message, s.cfg.OutboundCapacity), closed: make(chan struct{}), isDialer: true}
			s.addConn(h)
			s.cbConnect(h)
			backoff = s.cfg.ReconnectMin
			s.wg.Add(2)
			go s.reader(h)
			go s.writer(h)
			select {
			case <-h.closed:
			case <-s.ctx.Done():
				h.close()
				return
			}
		}
	}()
}

func setUnixNoSigPipe(_, _ string, c syscall.RawConn) error {
	var serr error
	_ = c.Control(func(fd uintptr) {
		const SO_NOSIGPIPE = 0x1022
		serr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, SO_NOSIGPIPE, 1)
	})
	return serr
}
