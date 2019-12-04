package watcher

import (
	"fmt"
	"os"
	"sync"
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
	wevent := Event{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 2; i++ {
			select {
			case wevent = <-watcher.GetEventChan():
				log.Info(wevent)
			case err := <-watcher.GetErrorChan():
				t.Error(err)
			case <-time.After(1 * time.Second):
				t.Error("No event was generated")
			}
		}
		fmt.Println("Done")
		wg.Done()
	}()

	file, err := os.Create("./sample_generated.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	file.Close()
	os.Remove("./sample_generated.yaml")
	// select {
	// case wevent = <-watcher.GetEventChan():
	// 	log.Info(wevent)
	// case err := <-watcher.GetErrorChan():
	// 	t.Error(err)
	// case <-time.After(1 * time.Second):
	// 	t.Error("No event was generated after removal of file")
	// }
	wg.Wait()
	watcher.Close()

	if f, ok := <-watcher.GetEventChan(); ok {
		t.Error("Channel is not closed")
		fmt.Println(f)
	}
}
