package watcher

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFsnotifyWrapper(t *testing.T) {
	watcher, err := NewFsWatcher()
	if err != nil {
		t.Error(err)
		return
	}
	err = watcher.Add("./")
	if err != nil {
		t.Error(err)
		return
	}
	file, err := os.Create("./sample_generated.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	file.Close()
	select {
	case wevent := <-watcher.GetEventChan():
		fmt.Println(wevent.Name)
	case err := <-watcher.GetErrorChan():
		t.Error(err)
	case <-time.After(1 * time.Second):
		t.Error("No event was generated after creation")
	}
	os.Remove("./sample_generated.yaml")
	select {
	case wevent := <-watcher.GetEventChan():
		fmt.Println(wevent.Name)
	case err := <-watcher.GetErrorChan():
		t.Error(err)
	case <-time.After(1 * time.Second):
		t.Error("No event was generated after removal of file")
	}
	watcher.Close()
	if _, ok := <-watcher.GetEventChan(); ok {
		t.Error("Channel is not closed")
	}
}
