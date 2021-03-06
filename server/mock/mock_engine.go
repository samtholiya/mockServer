package mock

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samtholiya/mockServer/server/model"
)

func (s *Server) getResponseForRequest(w http.ResponseWriter, r *http.Request, apis []model.API) {
	log.Tracef("%v headers", r.Header)
	for i := range apis {
		if s.compare.String(apis[i].Endpoint, r.URL.EscapedPath()) {
			log.Debugf("Url %v matched with %v", apis[i].Endpoint, r.URL.EscapedPath())
			scenarioNumber, scenario := s.getMatchedScenario(r, apis[i].Scenarios)
			s.writeResponse(w, scenario)
			if scenario.Delete {
				s.deleteScenario(r.Method, i, scenarioNumber)
			}
			return
		}
		log.Debugf("URL regex %v did not match %v", apis[i].Endpoint, r.URL.EscapedPath())
	}
	log.Warnf("No URL Matched %v", r.URL.EscapedPath())
	w.WriteHeader(404)
}

func (s *Server) deleteScenario(methodName string, apiNumber, scenarioNumber int) {
	s.app.API[methodName][apiNumber].Scenarios = append(s.app.API[methodName][apiNumber].Scenarios[:scenarioNumber], s.app.API[methodName][apiNumber].Scenarios[scenarioNumber+1:]...)
}

func (s *Server) writeResponse(w http.ResponseWriter, scenario model.Scenario) {

	if scenario.Response.Delay > 0 {
		time.Sleep(time.Duration(scenario.Response.Delay) * time.Second)
	}

	log.Debug("Writing response ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for key := range scenario.Response.Header {
		w.Header().Set(key, scenario.Response.Header[key][0])
	}
	if strings.Compare(scenario.Response.Payload.Type, "text") == 0 || strings.Compare(scenario.Response.Payload.Type, "json") == 0 {
		w.WriteHeader(scenario.Response.StatusCode)
		_, err := w.Write([]byte(scenario.Response.Payload.Data))
		if err != nil {
			log.Error(err)
		}
	}
	if strings.Compare(scenario.Response.Payload.Type, "base64") == 0 {
		w.WriteHeader(scenario.Response.StatusCode)
		data, err := base64.StdEncoding.DecodeString(scenario.Response.Payload.Data)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), 503)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			log.Error(err)
		}
	}
	if strings.Compare(scenario.Response.Payload.Type, "file") == 0 {
		file, err := os.Open(scenario.Response.Payload.Data)
		if err != nil {
			//File not found, send 404
			http.Error(w, "File not found.", 404)
			return
		}
		defer file.Close() //Close after function return

		//Get the Content-Type of the file
		//Create a buffer to store the header of the file in
		FileHeader := make([]byte, 512)
		//Copy the headers into the FileHeader buffer
		if _, err = file.Read(FileHeader); err != nil {
			log.Error(err)
		}
		//Get content type of file
		FileContentType := http.DetectContentType(FileHeader)

		//Get the file size
		FileStat, _ := file.Stat()                         //Get info from file
		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

		//Send the headers
		if len(w.Header().Get("Content-Type")) == 0 {
			w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
			w.Header().Set("Content-Type", FileContentType)
		}
		w.Header().Set("Content-Length", FileSize)
		w.WriteHeader(scenario.Response.StatusCode)
		//Send the file
		//We read 512 bytes from the file already, so we reset the offset back to 0
		if _, err = file.Seek(0, 0); err != nil {
			log.Error(err)
		} else {
			if _, err = io.Copy(w, file); err != nil {
				log.Error(err)
			}
		}
	}
}

func (s *Server) getMatchedScenario(r *http.Request, scenarios []model.Scenario) (int, model.Scenario) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	log.Tracef("Headers: %v", r.Header)
	for i := range scenarios {
		if !s.compare.MapStringArr(scenarios[i].Request.Header, r.Header) {
			log.Debugf("%v does not matches due to request headers", i)
			continue
		}
		if !s.compare.MapStringArr(scenarios[i].Request.Query, r.URL.Query()) {
			log.Debugf("%v does not matches due to request Query", i)
			continue
		}
		if scenarios[i].Request.Payload.Type == "text" {
			if !s.compare.String(scenarios[i].Request.Payload.Data, string(payload)) {
				log.Debugf("%v does not matches due to request Data", i)
				continue
			}
		} else if scenarios[i].Request.Payload.Type == "json" {
			if !s.compare.JSONString(scenarios[i].Request.Payload.Data, string(payload)) {
				log.Debugf("%v does not matches due to request data", i)
				continue
			}
		} else if scenarios[i].Request.Payload.Type == "base64" {
			if strings.Compare(scenarios[i].Request.Payload.Data, base64.StdEncoding.EncodeToString(payload)) != 0 {
				log.Debugf("%v does not matches due to request data", i)
				continue
			}
		} else if scenarios[i].Request.Payload.Type == "file" {
			data, err := ioutil.ReadFile(scenarios[i].Request.Payload.Data)
			if err != nil {
				log.Error(err)
			}
			if !bytes.Equal(payload, data) {
				log.Debugf("%v does not matches due to request data", i)
				continue
			}
		}
		log.Debugf("%v scenario matched", scenarios[i])
		return i, scenarios[i]
	}
	log.Debug("No scenario matched")
	return 0, model.Scenario{}
}
