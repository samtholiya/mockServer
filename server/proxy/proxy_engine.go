package proxy

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/samtholiya/mockServer/common"

	"github.com/samtholiya/mockServer/server/model"
	"gopkg.in/yaml.v2"
)

func (p *proxyServer) copyRequest(r *http.Request) (*http.Request, []byte) {
	body, _ := ioutil.ReadAll(r.Body)
	req, err := http.NewRequest(r.Method, p.Host+r.RequestURI, bytes.NewBuffer(body))
	if err != nil {
		p.log.Error(err)
		return nil, nil
	}
	req.Header = r.Header
	return req, body
}

func (p *proxyServer) copyToApp(r *http.Request, body []byte, resp *http.Response, respData []byte) {
	p.copyScenario(r, body, resp, respData)
	p.writeToFile()
}

func (p *proxyServer) copyScenario(r *http.Request, reqBody []byte, resp *http.Response, respData []byte) {
	scen := model.Scenario{}
	scen.Request.Header = r.Header
	scen.Request.Query = r.URL.Query()

	temp := &model.API{}
	flag := true
	for i := range p.app.API[r.Method] {
		if r.URL.Path == p.app.API[r.Method][i].Endpoint {
			temp = &p.app.API[r.Method][i]
			flag = false
			break
		}
	}
	if r.Header.Get("Content-Type") == "application/json" {
		scen.Request.Payload.Type = "json"
		scen.Request.Payload.Data = string(reqBody)
	}
	if strings.Contains(r.Header.Get("Content-Type"), "form") {
		p.log.Debug("Found Form encoded url")
		scen.Request.Payload.Type = "file"
		path := "./request_files/" + common.GetUniqueString(5) + ".req"
		os.MkdirAll("./request_files", os.ModePerm)
		tempFile, err := os.Create(path)
		if err != nil {
			p.log.Error(err)
			return
		}
		defer tempFile.Close()
		tempFile.Write(reqBody)
		if err != nil {
			panic(err)
		}
		scen.Request.Payload.Data = path
	}
	scen.Response = p.copyResponse(resp, respData)
	temp.Scenarios = append(temp.Scenarios, scen)
	if flag {
		temp.Endpoint = r.URL.Path
		p.app.API[r.Method] = append(p.app.API[r.Method], *temp)
	}
}

func (p proxyServer) writeToFile() {
	dataFound, _ := yaml.Marshal(p.app)
	file, err := os.Create(common.GetEnv("PROXY_GENERATED_CONFIG", "./proxy_generated.yaml"))
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
