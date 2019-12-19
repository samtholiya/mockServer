package proxy

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/samtholiya/mockServer/common"
	"github.com/samtholiya/mockServer/server/model"
	"github.com/samtholiya/mockServer/types"
	"github.com/sirupsen/logrus"
)

type proxyServer struct {
	Host       string
	app        *model.App
	log        *logrus.Logger
	client     *http.Client
	dataParser types.DataFormatParser
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
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(data); err != nil {
		p.log.Error(err)
	}
	go p.copyToApp(r, body, resp, data)
}

//NewProxyServer returns a proxy server
func NewProxyServer(host string, dataParser types.DataFormatParser) types.Proxy {
	app := model.New()
	proxy := proxyServer{
		log:        common.GetLogger(),
		Host:       host,
		app:        &app,
		dataParser: dataParser,
	}
	proxy.Init(true)
	return &proxy
}
