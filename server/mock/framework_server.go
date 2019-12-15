package mock

import (
	"net/http"

	"github.com/samtholiya/mockServer/server/model"

	"github.com/samtholiya/mockServer/types"
)

//Server contains watcher and other server
type Server struct {
	watch   types.Watcher
	compare types.Comparer
	app     *model.App
}

//SetWatcher sets the watcher for the server
func (s *Server) SetWatcher(watch types.Watcher) {
	s.watch = watch
}

//SetComparer sets the compare algo for the server
func (s *Server) SetComparer(compare types.Comparer) {
	s.compare = compare
}

//SetApp the App obj for the server
func (s *Server) SetApp(app model.App) {
	s.app = &app
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("%v Url %v method", r.URL, r.Method)
	s.handler(w, r)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	switch r.Method {
	case "GET":
		s.getResponseForRequest(w, r, s.app.Get)
	case "POST":
		s.getResponseForRequest(w, r, s.app.Post)
	}
}
