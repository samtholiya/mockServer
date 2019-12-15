package proxy

import (
	"net/http"

	"github.com/samtholiya/mockServer/server/model"
)

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
