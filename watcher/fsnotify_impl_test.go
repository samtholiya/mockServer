package watcher

import (
	"os"
	"testing"
	"time"

	"github.com/samtholiya/apiMocker/common"
)

func TestFsnotifyWrapper(t *testing.T) {
	log := common.GetLogger()
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
	var wevent Event

	file, err := os.Create("./sample_generated.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	file.Close()
	select {
	case wevent = <-watcher.GetEventChan():
		log.Info(wevent)
	case err := <-watcher.GetErrorChan():
		t.Error(err)
	case <-time.After(1 * time.Second):
		t.Error("No event was generated")
	}
	os.Remove("./sample_generated.yaml")
	select {
	case wevent = <-watcher.GetEventChan():
		log.Info(wevent)
	case err := <-watcher.GetErrorChan():
		t.Error(err)
	case <-time.After(1 * time.Second):
		t.Error("No event was generated after removal of file")
	}
	watcher.Close()

	if f, ok := <-watcher.GetEventChan(); ok {
		log.Error(f)
		t.Error("Channel is not closed")
	}
}
