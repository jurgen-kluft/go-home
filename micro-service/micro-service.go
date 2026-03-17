package microservice

import (
	"time"
)

// Delegate is a handler that the user can register for handling incoming messages
// from a specific peer path.
type Delegate func(m *Service, message *Message) bool

func (m *Service) RegisterHandler(peerPath string, delegate Delegate) {
	srcID := Hash32(peerPath)
	m.handlersByID[srcID] = delegate
	m.handlersByPath[peerPath] = delegate
}

func (m *Service) FindHandlerByID(srcID uint32) (delegate Delegate, exists bool) {
	if hander, exists := m.handlersByID[srcID]; exists {
		return hander, true
	}
	return nil, false
}

func (m *Service) FindHandlerByPath(peerPath string) (delegate Delegate, exists bool) {
	if hander, exists := m.handlersByPath[peerPath]; exists {
		return hander, true
	}
	return nil, false
}

func (m *Service) ConnectTo(peerPaths []string) {
	for _, r := range peerPaths {
		m.Connect(r)
	}
}

func (m *Service) Loop() {
	quit := false
	for !quit {
		connected := true
		ticker := time.NewTicker(m.tickFrequency)
		for connected {
			select {
			case <-ticker.C:
				if tickHandler, exists := m.FindHandlerByPath("tick"); exists {
					if !tickHandler(m, nil) {
						connected = false
						quit = true
					}
				}
			case msg := <-m.inboundCh:
				delegate, exists := m.FindHandlerByID(msg.SrcID())
				if exists {
					if !delegate(m, msg) {
						connected = false
						quit = true
					}
				} else {
					delegate, exists := m.FindHandlerByPath("*")
					if exists {
						if !delegate(m, msg) {
							connected = false
							quit = true
						}
					}

					if msg.SrcID() == Hash32("client/disconnected") {
						m.Logger.LogInfo("net", "disconnected")
						connected = false
					}
				}
			}
		}

		if !quit {
			m.Logger.LogInfo("net", "Waiting 5 seconds before re-connecting..")
			time.Sleep(5 * time.Second)
		} else {
			m.Logger.LogInfo("net", "End.")
		}
	}
}
