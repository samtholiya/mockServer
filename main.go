package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/samtholiya/apiMocker/comparer"

	"github.com/samtholiya/apiMocker/common"

	"github.com/samtholiya/apiMocker/watcher"

	"github.com/samtholiya/apiMocker/skeleton"
)

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	switch r.Method {
	case "GET":
		log.Println(r.URL)
	case "POST":

	}
	w.Write([]byte("Hello"))
}

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
	http.ListenAndServe(":3000", nil)
}

func SetupCloseHandler(server skeleton.Server) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Server")
		server.StopWatching()
		os.Exit(0)
	}()
}
