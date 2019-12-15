package proxy

import (
	"net/http"

	"github.com/samtholiya/mockServer/server/model"
)

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
