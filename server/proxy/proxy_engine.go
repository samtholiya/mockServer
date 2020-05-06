package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/samtholiya/mockServer/common"

	"github.com/samtholiya/mockServer/server/model"
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
		if err := os.MkdirAll("./request_files", os.ModePerm); err != nil {
			p.log.Error(err)
		}
		tempFile, err := os.Create(path)
		if err != nil {
			p.log.Error(err)
			return
		}
		defer tempFile.Close()
		if _, err = tempFile.Write(reqBody); err != nil {
			p.log.Error(err)
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
	if err := p.dataParser.WriteToFile(p.app, common.GetEnv("PROXY_GENERATED_CONFIG", "./proxy_generated.yaml")); err != nil {
		p.log.Error(err)
	}
}

func (p *proxyServer) copyResponse(resp *http.Response, respData []byte) model.Response {
	res := model.Response{}
	res.Header = resp.Header
	if len(respData) == 0 {
		res.Payload.Type = "text"
		res.Payload.Data = ""
		p.log.Trace("Empty response Data")
	} else if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		res.Payload.Type = "json"
		res.Payload.Data = string(respData)
	} else if strings.Contains(resp.Header.Get("Content-Type"), "text") {
		res.Payload.Type = "text"
		res.Payload.Data = string(respData)
	} else {
		res.Payload.Type = "file"
		common.CreateFolder("responseFiles")
		path := filepath.Join("responseFiles", common.GetUniqueString(5)+".res")
		f, err := os.Create(path)
		if err != nil {
			p.log.Error(err)
		}
		defer f.Close()
		_, err = f.Write(respData)
		if err != nil {
			p.log.Error(err)
		}
		res.Payload.Data = path
	}
	res.StatusCode = resp.StatusCode
	return res
}
