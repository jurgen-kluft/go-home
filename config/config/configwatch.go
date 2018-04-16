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

type configFileWatcher struct {
	configFiles map[string]*configFileWatchItem
}

func newConfigFileWatcher() *configFileWatcher {
	watcher := &configFileWatcher{}
	watcher.configFiles = map[string]*configFileWatchItem{}
	return watcher
}

func (c *configFileWatcher) watchConfigFile(path string, user string) {
	wi := new(configFileWatchItem)
	wi.path = path
	wi.user = user
	wi.lastEvent = NONE

	fi, err := os.Stat(path)
	if err == nil {
		wi.statInfo = fi
	} else if os.IsNotExist(err) {
		wi.lastEvent = NOEXIST
	} else if os.IsPermission(err) {
		wi.lastEvent = NOPERM
	} else {
		wi.lastEvent = INVALID
	}

	c.configFiles[path] = wi

	return
}

type configFileEvent struct {
	Path  string
	User  string
	Event int
}

func (c *configFileWatcher) update() []configFileEvent {
	events := []configFileEvent{}
	for fname, fitem := range c.configFiles {
		if fitem.update() {
			event := configFileEvent{Path: fname, User: fitem.user, Event: fitem.lastEvent}
			events = append(events, event)
		}
	}

	return events
}

type configFileWatchItem struct {
	path      string
	user      string
	statInfo  os.FileInfo
	lastEvent int
}

func (wi *configFileWatchItem) update() bool {
	fi, err := os.Stat(wi.path)
	if err != nil {
		if os.IsNotExist(err) {
			if wi.lastEvent == NOEXIST {
				return false
			} else if wi.lastEvent == DELETED {
				wi.lastEvent = NOEXIST
				return false
			} else {
				wi.lastEvent = DELETED
				return true
			}
		} else if os.IsPermission(err) {
			if wi.lastEvent == NOPERM {
				return false
			}
			wi.lastEvent = NOPERM
			return true
		} else {
			wi.lastEvent = INVALID
			return false
		}
	}

	if wi.lastEvent == NOEXIST {
		wi.lastEvent = CREATED
		wi.statInfo = fi
		return true
	} else if fi.ModTime().After(wi.statInfo.ModTime()) {
		wi.statInfo = fi
		switch wi.lastEvent {
		case NONE, CREATED, NOPERM, INVALID:
			wi.lastEvent = MODIFIED
		case DELETED, NOEXIST:
			wi.lastEvent = CREATED
		}
		return true
	} else if fi.Mode() != wi.statInfo.Mode() {
		wi.lastEvent = PERM
		wi.statInfo = fi
		return true
	}
	return false
}
