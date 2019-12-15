package proxy

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/samtholiya/mockServer/server/model"
	"gopkg.in/yaml.v2"
)

func (p *proxyServer) copyRequest(r *http.Request) (*http.Request, []byte) {

	switch r.Method {
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		req, err := http.NewRequest(r.Method, p.Host+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.log.Error(err)
			return nil, nil
		}
		req.Header = r.Header
		return req, body

	case "GET":
		req, err := http.NewRequest(r.Method, p.Host+r.RequestURI, r.Body)
		if err != nil {
			p.log.Error(err)
			return nil, nil
		}
		req.Header = r.Header
		return req, nil
	}

	return nil, nil
}

func (p *proxyServer) copyToApp(r *http.Request, body []byte, resp *http.Response, respData []byte) {
	switch r.Method {
	case "POST":
		p.copyPostScenario(r, body, resp, respData)
	case "GET":
		p.copyGetScenario(r, body, resp, respData)
	}
	p.writeToFile()
}

func (p proxyServer) writeToFile() {
	dataFound, _ := yaml.Marshal(p.app)
	file, err := os.Create("./proxy_generated.yaml")
	if err != nil {
		p.log.Error(err)
		return
	}
	_, err = file.WriteString(string(dataFound))
	if err != nil {
		p.log.Error(err)
		return
	}
	file.Close()
}

func (p *proxyServer) copyResponse(resp *http.Response, respData []byte) model.Response {
	res := model.Response{}
	res.Header = resp.Header
	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		res.Payload.Type = "json"
		res.Payload.Data = string(respData)
	} else if strings.Contains(resp.Header.Get("Content-Type"), "text") {
		res.Payload.Type = "text"
		res.Payload.Data = string(respData)
	} else {
		res.Payload.Type = "base64"
		res.Payload.Data = base64.StdEncoding.EncodeToString([]byte(respData))
	}
	res.StatusCode = resp.StatusCode
	return res
}
