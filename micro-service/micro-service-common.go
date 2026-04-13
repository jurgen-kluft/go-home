package microservice

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	logpkg "github.com/jurgen-kluft/go-home/logging"
)

const (
	Magic      uint32 = 0x554D5131 // 'UMQ1'
	Version    uint8  = 1
	HeaderLen         = 32
	MaxPayload        = 64 * 1024 // 64KiB; adjust as needed
	MaxFrame          = HeaderLen + MaxPayload
)

const (
	iMagic    = 0
	iVersion  = 4
	iMsgType  = 5
	iFlags    = 6
	iReserved = 7
	iLen      = 8
	iSrc      = 12
	iDst      = 16
	iSeq      = 20
	iCRC      = 28
	iPayload  = HeaderLen
)

type Message struct {
	Frame   []byte
	Payload []byte
	pooled  bool
}

func (m *Message) MsgType() uint8      { return m.Frame[iMsgType] }
func (m *Message) Flags() uint8        { return m.Frame[iFlags] }
func (m *Message) Len() uint32         { return binary.LittleEndian.Uint32(m.Frame[iLen:]) }
func (m *Message) SrcID() uint32       { return binary.LittleEndian.Uint32(m.Frame[iSrc:]) }
func (m *Message) DstID() uint32       { return binary.LittleEndian.Uint32(m.Frame[iDst:]) }
func (m *Message) Seq() uint64         { return binary.LittleEndian.Uint64(m.Frame[iSeq:]) }
func (m *Message) HeaderBytes() []byte { return m.Frame[:HeaderLen] }

func (m *Message) SetMsgType(t uint8) { m.Frame[iMsgType] = t }
func (m *Message) SetFlags(f uint8)   { m.Frame[iFlags] = f }
func (m *Message) SetLen(n uint32)    { binary.LittleEndian.PutUint32(m.Frame[iLen:], n) }
func (m *Message) SetSrcID(id uint32) { binary.LittleEndian.PutUint32(m.Frame[iSrc:], id) }
func (m *Message) SetDstID(id uint32) { binary.LittleEndian.PutUint32(m.Frame[iDst:], id) }
func (m *Message) SetSeq(seq uint64)  { binary.LittleEndian.PutUint64(m.Frame[iSeq:], seq) }

func (m *Message) Release() {
	if m == nil {
		return
	}
	if m.pooled && cap(m.Frame) >= MaxFrame {
		framePool.Put(m.Frame[:MaxFrame])
	}
	m.Frame = nil
	m.Payload = nil
	msgPool.Put(m)
}

var (
	framePool = sync.Pool{New: func() any { return make([]byte, MaxFrame) }}
	msgPool   = sync.Pool{New: func() any { return &Message{} }}
)

func newMessage(payloadLen int) (*Message, error) {
	if payloadLen < 0 || payloadLen > MaxPayload {
		return nil, fmt.Errorf("payload too large: %d", payloadLen)
	}
	m := msgPool.Get().(*Message)
	frame := framePool.Get().([]byte)
	m.Frame = frame[:HeaderLen+payloadLen]
	m.Payload = m.Frame[HeaderLen : HeaderLen+payloadLen]
	m.pooled = true
	binary.LittleEndian.PutUint32(m.Frame[iMagic:], Magic)
	m.Frame[iVersion] = Version
	return m, nil
}

func cloneMessage(src *Message) *Message {
	n := int(src.Len())
	m, _ := newMessage(n)
	copy(m.Frame[:HeaderLen+n], src.Frame[:HeaderLen+n])
	return m
}

type Config struct {
	ListenPath       string
	Peers            []string
	InboundCapacity  int
	OutboundCapacity int
	ReconnectMin     time.Duration
	ReconnectMax     time.Duration
	DropWhenFull     bool

	// Optional observability callbacks
	OnConnect    func(peerPath string, peerID uint32)
	OnDisconnect func(peerPath string, peerID uint32)
}

type Service struct {
	Name string

	Logger *logpkg.Logger

	handlersByID   map[uint32]Delegate
	handlersByPath map[string]Delegate

	cfg       Config
	inboundCh chan *Message
	srcID     uint32

	ln     *net.UnixListener
	connMu sync.RWMutex
	conns  map[string]*connHandle

	ctx    context.Context
	cancel context.CancelFunc

	seq uint64
	wg  sync.WaitGroup

	tickMessage   *Message
	tickFrequency time.Duration
}

type connHandle struct {
	parent    *Service
	peerPath  string
	peerID    uint32
	conn      *net.UnixConn
	writeCh   chan *Message
	closed    chan struct{}
	closeOnce sync.Once
	isDialer  bool
}

func Hash32(p string) uint32 {
	h := sha1.Sum([]byte(p))
	return binary.LittleEndian.Uint32(h[:4])
}

func New(name string, peerPath string, tickFrequency time.Duration) (*Service, error) {
	cfg := Config{ListenPath: peerPath}

	if cfg.ListenPath == "" {
		return nil, errors.New("ListenPath required")
	}
	if cfg.InboundCapacity <= 0 {
		cfg.InboundCapacity = 1024
	}
	if cfg.OutboundCapacity <= 0 {
		cfg.OutboundCapacity = 1024
	}
	if cfg.ReconnectMin <= 0 {
		cfg.ReconnectMin = 100 * time.Millisecond
	}
	if cfg.ReconnectMax <= 0 {
		cfg.ReconnectMax = 5 * time.Second
	}

	_ = os.MkdirAll(filepath.Dir(cfg.ListenPath), 0o755)
	if err := os.RemoveAll(cfg.ListenPath); err != nil {
		return nil, fmt.Errorf("remove stale socket: %w", err)
	}

	addr := &net.UnixAddr{Name: cfg.ListenPath, Net: "unix"}
	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		return nil, fmt.Errorf("listen unix: %w", err)
	}
	_ = os.Chmod(cfg.ListenPath, 0o600)

	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		Name:           name,
		Logger:         logpkg.New(name),
		handlersByID:   make(map[uint32]Delegate),
		handlersByPath: make(map[string]Delegate),
		cfg:            cfg,
		inboundCh:      make(chan *Message, cfg.InboundCapacity),
		srcID:          Hash32(cfg.ListenPath),
		ln:             ln,
		conns:          make(map[string]*connHandle),
		ctx:            ctx,
		cancel:         cancel,
		tickFrequency:  tickFrequency,
	}

	s.wg.Add(1)
	go s.acceptLoop()
	for _, p := range cfg.Peers {
		if p == cfg.ListenPath {
			continue
		}
		s.startDialer(p)
	}

	s.tickMessage = &Message{}
	s.tickMessage.SetMsgType(0) // Tick message type
	s.tickMessage.SetSrcID(Hash32("tick"))
	s.tickMessage.SetLen(0)

	go func() {
		time.Sleep(s.tickFrequency)
		s.inboundCh <- s.tickMessage
	}()

	return s, nil
}

// Connect starts (or restarts) a client-side connector (dialer) to the given server socket path.
// It returns immediately and reconnects automatically on failures.
func (s *Service) Connect(peerPath string) {
	if peerPath == "" || peerPath == s.cfg.ListenPath {
		return
	}
	if s.getConn(peerPath) != nil {
		return
	}
	s.startDialer(peerPath)
}

// ConnectAndWait connects to peerPath and waits until the connection is established or timeout elapses.
// Returns true if connected within the timeout, false otherwise.
func (s *Service) ConnectAndWait(peerPath string, timeout time.Duration) bool {
	s.Connect(peerPath)
	return s.ReadyWait(peerPath, timeout)
}

// Ready returns connectivity for the provided peers. If no peers are given, it returns
// readiness for all known connections.
func (s *Service) Ready(peers ...string) map[string]bool {
	res := make(map[string]bool)
	if len(peers) == 0 {
		s.connMu.RLock()
		for p, h := range s.conns {
			res[p] = (h != nil && !isClosed(h))
		}
		s.connMu.RUnlock()
		return res
	}
	for _, p := range peers {
		c := s.getConn(p)
		res[p] = (c != nil && !isClosed(c))
	}
	return res
}

// ReadyWait blocks until the connection to peerPath is up, or timeout expires.
// Returns true if the connection is up before the deadline.
func (s *Service) ReadyWait(peerPath string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		c := s.getConn(peerPath)
		if c != nil && !isClosed(c) {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		select {
		case <-time.After(20 * time.Millisecond):
		case <-s.ctx.Done():
			return false
		}
	}
}

// Peers returns a sorted list of peer paths this service currently knows about (connected or in reconnect),
// filtered to only those with an active handle present.
func (s *Service) Peers() []string {
	s.connMu.RLock()
	peers := make([]string, 0, len(s.conns))
	for p, h := range s.conns {
		if h != nil && !isClosed(h) {
			peers = append(peers, p)
		}
	}
	s.connMu.RUnlock()
	sort.Strings(peers)
	return peers
}

// Disconnect intentionally drops a client connection to the given peer path.
// If no such connection exists, it's a no-op.
func (s *Service) Disconnect(peerPath string) {
	c := s.getConn(peerPath)
	if c != nil {
		c.close()
	}
}

func (s *Service) Close() error {
	s.cancel()
	_ = s.ln.Close()
	s.connMu.RLock()
	for _, c := range s.conns {
		c.close()
	}
	s.connMu.RUnlock()
	s.wg.Wait()
	_ = os.RemoveAll(s.cfg.ListenPath)
	close(s.inboundCh)
	return nil
}

func (s *Service) Inbound() <-chan *Message { return s.inboundCh }

func (s *Service) NewMessage(msgType uint8, n int) (*Message, error) {
	m, err := newMessage(n)
	if err != nil {
		return nil, err
	}
	m.SetMsgType(msgType)
	m.SetLen(uint32(n))
	m.SetSrcID(s.srcID)
	m.SetSeq(atomic.AddUint64(&s.seq, 1))
	return m, nil
}

func (s *Service) NewTextMessage(text string) (*Message, error) {
	m, err := s.NewMessage(1, len(text))
	if err != nil {
		return nil, err
	}
	copy(m.Payload, []byte(text))
	return m, nil
}

func (s *Service) SendJsonTo(peerPath string, json string) error {
	msg, err := s.NewTextMessage(json)
	if err != nil {
		return err
	}
	return s.SendTo(peerPath, msg)
}

func (s *Service) SendTo(peerPath string, m *Message) error {
	c := s.getConn(peerPath)
	if c == nil {
		return fmt.Errorf("no connection to %s yet", peerPath)
	}
	m.SetDstID(c.peerID)
	toSend := m
	select {
	case c.writeCh <- toSend:
		return nil
	default:
		if s.cfg.DropWhenFull {
			return errors.New("outbound channel full (dropped)")
		}
		select {
		case c.writeCh <- toSend:
			return nil
		case <-s.ctx.Done():
			return context.Canceled
		}
	}
}

func (s *Service) Broadcast(m *Message) int {
	sent := 0
	for _, c := range s.conns {
		if c == nil {
			continue
		}
		cp := cloneMessage(m)
		cp.SetDstID(c.peerID)
		select {
		case c.writeCh <- cp:
			sent++
		default:
			if !s.cfg.DropWhenFull {
				select {
				case c.writeCh <- cp:
					sent++
				case <-s.ctx.Done():
					cp.Release()
					return sent
				}
			} else {
				cp.Release()
			}
		}
	}
	return sent
}

func (s *Service) getConn(peerPath string) *connHandle {
	s.connMu.RLock()
	c := s.conns[peerPath]
	s.connMu.RUnlock()
	return c
}

func (s *Service) addConn(c *connHandle) {
	s.connMu.Lock()
	s.conns[c.peerPath] = c
	s.connMu.Unlock()
}

func (s *Service) delConn(peerPath string) {
	s.connMu.Lock()
	delete(s.conns, peerPath)
	s.connMu.Unlock()
}

func (s *Service) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.ln.AcceptUnix()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				time.Sleep(50 * time.Millisecond)
				continue
			}
		}
		_ = conn.SetReadBuffer(MaxFrame * 2)
		_ = conn.SetWriteBuffer(MaxFrame * 2)
		h := &connHandle{parent: s, peerPath: "<accepted>", peerID: 0, conn: conn, writeCh: make(chan *Message, s.cfg.OutboundCapacity), closed: make(chan struct{})}
		s.addConn(h)
		s.cbConnect(h)
		s.wg.Add(2)
		go s.reader(h)
		go s.writer(h)
	}
}

func (h *connHandle) close() {
	h.closeOnce.Do(func() {
		close(h.closed)
		_ = h.conn.Close()
		close(h.writeCh)
		if h.parent != nil {
			h.parent.delConn(h.peerPath)
			h.parent.cbDisconnect(h)
		}
	})
}

func (s *Service) reader(h *connHandle) {
	defer s.wg.Done()
	defer h.close()
	header := make([]byte, HeaderLen)
	for {
		if err := readFull(h.conn, header); err != nil {
			return
		}
		magic := binary.LittleEndian.Uint32(header[iMagic:])
		if magic != Magic || header[iVersion] != Version {
			return
		}
		length := binary.LittleEndian.Uint32(header[iLen:])
		if length > MaxPayload {
			return
		}
		m, _ := newMessage(int(length))
		copy(m.Frame[:HeaderLen], header)
		if length > 0 {
			if err := readFull(h.conn, m.Payload); err != nil {
				m.Release()
				return
			}
		}
		if h.peerID == 0 {
			h.peerID = m.SrcID()
		}
		select {
		case s.inboundCh <- m:
		case <-s.ctx.Done():
			m.Release()
			return
		}
	}
}

func (s *Service) writer(h *connHandle) {
	defer s.wg.Done()
	defer h.close()
	for {
		var m *Message
		var ok bool
		select {
		case m, ok = <-h.writeCh:
			if !ok {
				return
			}
		case <-s.ctx.Done():
			return
		}
		if m.DstID() == 0 {
			m.SetDstID(h.peerID)
		}
		if err := writeFull(h.conn, m.Frame[:HeaderLen+int(m.Len())]); err != nil {
			return
		}
	}
}

// --- Observability helpers ---

func (s *Service) cbConnect(h *connHandle) {
	if s.cfg.OnConnect != nil {
		go s.cfg.OnConnect(h.peerPath, h.peerID)
	}
}

func (s *Service) cbDisconnect(h *connHandle) {
	if s.cfg.OnDisconnect != nil {
		go s.cfg.OnDisconnect(h.peerPath, h.peerID)
	}
}

// --- Utilities ---

func readFull(r io.Reader, buf []byte) error { _, err := io.ReadFull(r, buf); return err }

func writeFull(w io.Writer, buf []byte) error {
	for len(buf) > 0 {
		n, err := w.Write(buf)
		if n > 0 {
			buf = buf[n:]
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func isClosed(h *connHandle) bool {
	select {
	case <-h.closed:
		return true
	default:
		return false
	}
}

func jitter(d time.Duration) time.Duration {
	n := time.Duration(int64(d) + (int64(d)/5)*(int64(time.Now().UnixNano())%2*2-1))
	if n < 0 {
		return d
	}
	return n
}
