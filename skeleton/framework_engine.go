package skeleton

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (s *Server) getResponseForRequest(w http.ResponseWriter, r *http.Request, apis []API) {
	for i := range apis {
		if s.compare.String(apis[i].Endpoint, r.URL.String()) {
			scenario := s.getMatchedScenario(r, apis[i].Scenarios)
			s.writeResponse(w, scenario)
			return
		}
	}
}

func (s *Server) writeResponse(w http.ResponseWriter, scenario Scenario) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for key := range scenario.Response.Header {
		w.Header().Set(key, scenario.Response.Header[key][0])
	}
	if strings.Compare(scenario.Response.Payload.Type, "text") == 0 {
		_, err := w.Write([]byte(scenario.Response.Payload.Data))
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
		w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
		w.Header().Set("Content-Type", FileContentType)
		w.Header().Set("Content-Length", FileSize)

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
	w.WriteHeader(scenario.Response.StatusCode)
}

func (s *Server) getMatchedScenario(r *http.Request, scenarios []Scenario) Scenario {
	for i := range scenarios {
		if !s.compare.MapStringArr(scenarios[i].Request.Header, r.Header) {
			continue
		}
		if !s.compare.MapStringArr(scenarios[i].Request.Query, r.URL.Query()) {
			continue
		}
		if scenarios[i].Request.Payload.Type == "text" {
			payload, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error(err)
				continue
			}
			if !s.compare.String(scenarios[i].Request.Payload.Data, string(payload)) {
				continue
			}
		} else if scenarios[i].Request.Payload.Type == "json" {
			payload, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error(err)
				continue
			}
			if !s.compare.String(scenarios[i].Request.Payload.Data, string(payload)) {
				continue
			}
		}
		return scenarios[i]
	}
	return Scenario{}
}
