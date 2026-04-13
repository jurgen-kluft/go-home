//go:build darwin

package filewatch

import (
    "errors"
    "os"
    "sync"
    "syscall"
)

type Callback func(path string, flags uint32)

type Watcher struct {
    kq     int
    mu     sync.Mutex
    files  map[int]*watchedFile
    closed bool
}

type watchedFile struct {
    path string
    file *os.File
    cb   Callback
}

// NewWatcher creates a kqueue-backed file watcher.
func NewWatcher() (*Watcher, error) {
    kq, err := syscall.Kqueue()
    if err != nil {
        return nil, err
    }

    return &Watcher{
        kq:    kq,
        files: make(map[int]*watchedFile),
    }, nil
}

// Add starts watching a specific file path.
func (w *Watcher) Add(path string, cb Callback) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.closed {
        return errors.New("watcher closed")
    }

    f, err := os.Open(path)
    if err != nil {
        return err
    }

    fd := int(f.Fd())

    ev := syscall.Kevent_t{
