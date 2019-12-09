package skeleton

import (
	"io/ioutil"
	"runtime"

	"github.com/samtholiya/apiMocker/types"

	"github.com/samtholiya/apiMocker/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var log *logrus.Logger

func init() {
	log = common.GetLogger()
}

//GetApp Returns currently loaded App object
func (s *Server) GetApp() *App {
	return s.app
}

func (s *Server) appLoader(file string) {
	s.loadDataFromFile(file)
	for {
		select {
		case event, ok := <-s.watch.GetEventChan():
			if !ok {
				log.Info("Exiting goroutine for monitoring file change")
				return
			}
			if event.Operation&types.Write == types.Write {
				s.loadDataFromFile(file)
			} else {
				log.Info("Operation was not write")
			}
		case err := <-s.watch.GetErrorChan():
			log.Error(err)
		}
	}
}

func (s *Server) loadDataFromFile(file string) {
	dataBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("Error reading file %v", file)
	} else {
		err = yaml.Unmarshal(dataBytes, s.app)
		if err != nil {
			log.Error(err)
		}
		log.Infof("New app structure %v", s.app)
	}
}

//StartWatching starts watching on config file and loads eventually
func (s *Server) StartWatching() {
	file := common.GetEnv("MOCK_CONFIG", "./config.yaml")
	if err := s.watch.Add(file); err != nil {
		log.Error(err)
		return
	}
	go s.appLoader(file)
	runtime.Gosched()
}

//StopWatching stops watching config file
func (s *Server) StopWatching() {
	s.watch.Close()
}
