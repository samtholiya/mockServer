package skeleton

import (
	"io/ioutil"

	"github.com/samtholiya/apiMocker/common"
	"github.com/samtholiya/apiMocker/watcher"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log *logrus.Logger

func init() {
	log = common.GetLogger()
}

//GetApp Returns currently loaded App object
func (s *Server) GetApp() App {
	return s.app
}

func (s *Server) appLoader(file string) {
	s.app = loadDataFromFile(file)
	for {
		select {
		case event, ok := <-s.watch.GetEventChan():
			if !ok {
				log.Info("Exiting goroutine for monitoring file change")
				return
			}
			if event.Operation&watcher.Write != 0 {
				s.app = loadDataFromFile(file)
			} else {
				log.Info("Operation was not write")
			}
		case err := <-s.watch.GetErrorChan():
			log.Error(err)
		}
	}
}

func loadDataFromFile(file string) App {
	app := App{}
	dataBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("Error reading file %v", file)
	} else {

		err = yaml.Unmarshal(dataBytes, &app)
		if err != nil {
			log.Error(err)
		}
		log.Infof("New app structure %v", app)
	}
	return app
}

//StartWatching starts watching on config file and loads eventually
func (s Server) StartWatching() {
	file := common.GetEnv("MOCK_CONFIG", "./config.yaml")
	if err := s.watch.Add(file); err != nil {
		log.Error(err)
		return
	}
	go s.appLoader(file)
}

//StopWatching stops watching config file
func (s Server) StopWatching() {
	s.watch.Close()
}
