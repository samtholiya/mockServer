package watcher

import (
	"context"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/samtholiya/apiMocker/common"
	"github.com/sirupsen/logrus"
)

type fsnoti struct {
	*fsnotify.Watcher
	ctx          context.Context
	cancelCtx    context.CancelFunc
	log          *logrus.Logger
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
	common.GetLogger()
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	fs := &fsnoti{
		Watcher:      watcher,
		ctx:          ctx,
		cancelCtx:    cancelCtx,
		log:          common.GetLogger(),
		eventChannel: eventChan,
	}
	go fs.runAsyncEventConv()

	return fs, nil
}

func (fs *fsnoti) runAsyncEventConv() {
	for {
		eventValue, ok := <-fs.Watcher.Events

		if ok {
			temp := fs.parseEvent(eventValue)
			select {
			case fs.eventChannel <- temp:
				fs.log.Tracef("Sent %v to eventChannel", temp)
			case <-fs.ctx.Done():
				close(fs.eventChannel)
				fs.log.Trace("Closed eventChannel")
				return
			}
		} else {
			close(fs.eventChannel)
			fs.log.Trace("Closed eventChannel")
			return
		}
	}

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
	fs.cancelCtx()
	runtime.Gosched()
	return nil
}
