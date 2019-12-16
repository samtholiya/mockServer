package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/samtholiya/mockServer/server/proxy"

	"github.com/samtholiya/mockServer/comparer"
	"github.com/samtholiya/mockServer/server/mock"
	"github.com/samtholiya/mockServer/server/model"
	"github.com/samtholiya/mockServer/watcher"
	"github.com/sirupsen/logrus"

	"github.com/samtholiya/mockServer/common"
)

var log *logrus.Logger

func main() {
	isProxyServer := flag.Bool("proxy", false, "Run in proxy server mode")
	host := flag.String("host", "https://httpbin.org", "Server url for proxy server")
	port := flag.String("port", "3000", "Port number for the server")
	flag.Parse()
	if *isProxyServer {
		startProxyServer(*host, *port)
	}
	startMockServer(*port)
}

func startMockServer(port string) {
	log = common.GetLogger()
	server := mock.Server{}
	watch, err := watcher.NewFsWatcher()
	if err != nil {
		log.Fatal(err)
	}
	server.SetApp(model.App{})
	server.SetWatcher(watch)
	server.SetComparer(comparer.NewRegexComparer())
	go SetupCloseHandler(server)
	server.StartWatching()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})
	if err = http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err)
	}

}

func startProxyServer(host, port string) {
	pServer := proxy.NewProxyServer(host)
	http.Handle("/", pServer)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err)
	}
}

func SetupCloseHandler(server mock.Server) {
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
