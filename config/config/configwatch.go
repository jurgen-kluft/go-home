package main

import (
	"os"
)

// These values represent the events fswatch knows about. Watch uses a
// fsstat call to look up file information; a file will only have a NOPERM
// event if the parent directory has no search permission (i.e. parent
// directory doesn't have executable permissions for the current user).
const (
	NONE     = iota // No event, initial state.
	CREATED         // File was created.
	DELETED         // File was deleted.
	MODIFIED        // File was modified.
	PERM            // Changed permissions
	NOEXIST         // File does not exist.
	NOPERM          // No permissions for the file (see const block comment).
	INVALID         // Any type of error not represented above.
)

type ConfigFileWatcher struct {
	ConfigFiles map[string]*ConfigFileWatchItem
}

func NewConfigFileWatcher() *ConfigFileWatcher {
	watcher := &ConfigFileWatcher{}
	watcher.ConfigFiles = map[string]*ConfigFileWatchItem{}
	return watcher
}

func (c *ConfigFileWatcher) WatchConfigFile(path string, user string) {
	wi := new(ConfigFileWatchItem)
	wi.Path = path
	wi.User = user
	wi.LastEvent = NONE

	fi, err := os.Stat(path)
	if err == nil {
		wi.StatInfo = fi
	} else if os.IsNotExist(err) {
		wi.LastEvent = NOEXIST
	} else if os.IsPermission(err) {
		wi.LastEvent = NOPERM
	} else {
		wi.LastEvent = INVALID
	}

	c.ConfigFiles[path] = wi

	return
}

type ConfigFileEvent struct {
	Path  string
	User  string
	Event int
}

func (w *ConfigFileWatcher) Update() []ConfigFileEvent {
	events := []ConfigFileEvent{}

	for fname, fitem := range w.ConfigFiles {
		if fitem.Update() {
			event := ConfigFileEvent{Path: fname, User: fitem.User, Event: fitem.LastEvent}
			events = append(events, event)
		}
	}

	return events
}

type ConfigFileWatchItem struct {
	Path      string
	User      string
	StatInfo  os.FileInfo
	LastEvent int
}

func (wi *ConfigFileWatchItem) Update() bool {
	fi, err := os.Stat(wi.Path)
	if err != nil {
		if os.IsNotExist(err) {
			if wi.LastEvent == NOEXIST {
				return false
			} else if wi.LastEvent == DELETED {
				wi.LastEvent = NOEXIST
				return false
			} else {
				wi.LastEvent = DELETED
				return true
			}
		} else if os.IsPermission(err) {
			if wi.LastEvent == NOPERM {
				return false
			} else {
				wi.LastEvent = NOPERM
				return true
			}
		} else {
			wi.LastEvent = INVALID
			return false
		}
	}

	if wi.LastEvent == NOEXIST {
		wi.LastEvent = CREATED
		wi.StatInfo = fi
		return true
	} else if fi.ModTime().After(wi.StatInfo.ModTime()) {
		wi.StatInfo = fi
		switch wi.LastEvent {
		case NONE, CREATED, NOPERM, INVALID:
			wi.LastEvent = MODIFIED
		case DELETED, NOEXIST:
			wi.LastEvent = CREATED
		}
		return true
	} else if fi.Mode() != wi.StatInfo.Mode() {
		wi.LastEvent = PERM
		wi.StatInfo = fi
		return true
	}
	return false
}
