package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/samtholiya/mockServer/comparer"

	"github.com/samtholiya/mockServer/common"

	"github.com/samtholiya/mockServer/watcher"

	"github.com/samtholiya/mockServer/skeleton"
)

func main() {
	// finalHandler := http.HandlerFunc(final)
	// watch, err := watcher.NewFsWatcher()
	// if err != nil {
	// 	log.Println(err)
	// 	log.Fatal(err)
	// 	return
	// }
	// server := skeleton.Server{}
	// server.SetWatcher(watch)
	// server.StartWatching()
	// http.Handle("/", (finalHandler))
	// http.ListenAndServe(":3000", nil)
	log := common.GetLogger()
	server := skeleton.Server{}
	watch, err := watcher.NewFsWatcher()
	if err != nil {
		log.Fatal(err)
	}
	server.SetApp(skeleton.App{})
	server.SetWatcher(watch)
	server.SetComparer(comparer.NewRegexComparer())
	go SetupCloseHandler(server)
	server.StartWatching()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})
	if err = http.ListenAndServe(":3000", nil); err != nil {
		log.Error(err)
	}
}

func SetupCloseHandler(server skeleton.Server) {
	c := make(chan os.Signal, 2)
	log.Println("Press Ctrl+C to close the server")
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Server")
		server.StopWatching()
		os.Exit(0)
	}()
}
