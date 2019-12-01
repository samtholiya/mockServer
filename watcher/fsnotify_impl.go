package watcher

import (
	"github.com/fsnotify/fsnotify"
)

type fsnoti struct {
	*fsnotify.Watcher
	eventChannel chan Event
}

//GetErrorChan returns channel which will output errors in the watcher
func (fs *fsnoti) GetErrorChan() chan error {
	return fs.Errors
}

//NewFsWatcher returns a watcher with implementation by fsnotify library
func NewFsWatcher() (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	eventChan := make(chan Event)

	fs := &fsnoti{
		Watcher:      watcher,
		eventChannel: eventChan,
	}

	go func(eventChannel chan Event, fsEventChannel chan fsnotify.Event, parseEvent func(fsnotify.Event) Event) {
		for {
			eventValue, ok := <-fsEventChannel
			if ok {
				eventChannel <- parseEvent(eventValue)
				continue
			}
		}
	}(eventChan, fs.Events, fs.parseEvent)

	return fs, nil
}

//GetEventChan returns channel which will output file system notification.
func (fs *fsnoti) GetEventChan() chan Event {
	return fs.eventChannel
}

func (fs fsnoti) parseEvent(event fsnotify.Event) Event {
	tempValue := Event{}
	tempValue.Name = event.Name
	switch event.Op {
	case fsnotify.Create:
		tempValue.Operation = Create
	case fsnotify.Remove:
		tempValue.Operation = Remove
	case fsnotify.Write:
		tempValue.Operation = Write
	case fsnotify.Rename:
		tempValue.Operation = Rename
	case fsnotify.Chmod:
		tempValue.Operation = Chmod
	}
	return tempValue
}


func (fs *fsnoti) Close() error {
	if err := fs.Watcher.Close(); err != nil {
		return err
	}
	close(fs.eventChannel)
	return nil
}
