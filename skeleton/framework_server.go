package skeleton

import (
	"github.com/samtholiya/apiMocker/watcher"
)

//Server contains watcher and other server
type Server struct {
	watch watcher.Watcher
	app   App
}

//SetWatcher sets the watcher for the server
func (s *Server) SetWatcher(watch watcher.Watcher) {
	s.watch = watch
}
