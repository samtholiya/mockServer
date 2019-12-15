package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/samtholiya/mockServer/common"
	"github.com/samtholiya/mockServer/server/model"
	"github.com/samtholiya/mockServer/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type proxyServer struct {
	Host   string
	app    *model.App
	log    *logrus.Logger
	client *http.Client
}

func (p *proxyServer) Init(insecure bool) {
	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		p.client = &http.Client{Transport: tr}
	} else {
		p.client = &http.Client{}
	}
}

func (p *proxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.log.Infof("Rediriecting %v to %v", r.URL, p.Host)
	req, body := p.copyRequest(r)
	resp, err := p.client.Do(req)
	if err != nil {
		p.log.Error(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.log.Error(err)
	}
	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(data); err != nil {
		p.log.Error(err)
	}
	go p.copyToApp(r, body, resp, data)
}

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
	//	p.log.Println(p.app)
	p.writeToFile()
}

func (p *proxyServer) copyGetScenario(r *http.Request, body []byte, resp *http.Response, respData []byte) {
	scen := model.Scenario{}
	scen.Request.Header = r.Header
	scen.Request.Query = r.URL.Query()
	temp := &model.API{}
	flag := true
	for i := range p.app.Get {
		if r.URL.Path == p.app.Get[i].Endpoint {
			temp = &p.app.Get[i]
			flag = false
			break
		}
	}
	scen.Response = p.copyResponse(resp, respData)
	temp.Scenarios = append(temp.Scenarios, scen)
	if flag {
		temp.Endpoint = r.URL.Path
		p.app.Get = append(p.app.Get, *temp)
	}

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

func (p *proxyServer) copyPostScenario(r *http.Request, body []byte, resp *http.Response, respData []byte) {
	scen := model.Scenario{}
	scen.Request.Header = r.Header
	scen.Request.Query = r.URL.Query()

	temp := &model.API{}
	flag := true
	for i := range p.app.Post {
		if r.URL.Path == p.app.Post[i].Endpoint {
			temp = &p.app.Post[i]
			flag = false
			break
		}
	}
	if r.Header.Get("Content-Type") == "application/json" {
		scen.Request.Payload.Type = "json"
		scen.Request.Payload.Data = string(body)
	}
	scen.Response = p.copyResponse(resp, respData)
	temp.Scenarios = append(temp.Scenarios, scen)
	if flag {
		temp.Endpoint = r.URL.Path
		p.app.Post = append(p.app.Post, *temp)
	}
}

//NewProxyServer returns a proxy server
func NewProxyServer(host string) types.Proxy {
	proxy := proxyServer{
		log:  common.GetLogger(),
		Host: host,
		app:  &model.App{},
	}
	proxy.Init(true)
	return &proxy
}
