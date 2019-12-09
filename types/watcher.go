package types

import "bytes"

// Watcher watches a set of files, delivering events to a channel.
type Watcher interface {

	// Add starts watching the named file or directory (non-recursively).
	Add(fileFolder string) error

	// Remove stops watching the the named file or directory (non-recursively).
	Remove(fileFolder string) error

	// Close removes all watches and closes the events channel.
	Close() error

	//GetErrorChan returns channel which will output errors in the watcher
	GetErrorChan() chan error

	//GetEventChan returns channel which will output file system notification.
	GetEventChan() chan Event
}

// Event represents a single file system notification.
type Event struct {
	Name      string    // Relative path to the file or directory.
	Operation Operation // File operation that triggered the event.
}

// Operation describes a set of file operations.
type Operation uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Operation = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

func (op Operation) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if op&Create == Create {
		buffer.WriteString("|CREATE")
	}
	if op&Remove == Remove {
		buffer.WriteString("|REMOVE")
	}
	if op&Write == Write {
		buffer.WriteString("|WRITE")
	}
	if op&Rename == Rename {
		buffer.WriteString("|RENAME")
	}
	if op&Chmod == Chmod {
		buffer.WriteString("|CHMOD")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:] // Strip leading pipe
}
