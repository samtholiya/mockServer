package mock

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/samtholiya/mockServer/comparer"
	"github.com/samtholiya/mockServer/server/model"
	"github.com/samtholiya/mockServer/types"
)

func TestServerHTTPGet(t *testing.T) {
	log.SetLevel(logrus.DebugLevel)
	server := Server{}
	server.SetWatcher(types.TestWatcher{})
	server.SetComparer(comparer.NewRegexComparer())
	app := model.App{
		API: map[string][]model.API{
			"GET": []model.API{model.API{
				Endpoint:    "/get",
				Description: "Here is the description",
				Scenarios: []model.Scenario{
					model.Scenario{
						Request: model.Request{
							Header: map[string][]string{"Accept": []string{"application/json"}},
						},
						Response: model.Response{
							Payload: model.Payload{
								Type: "text",
								Data: "Hello World",
							},
						},
					},
				},
			},
			},
		},
	}
	server.SetApp(app)
	req := httptest.NewRequest("GET", "/get", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if (resp.StatusCode) != 200 {
		t.Error("Status code should be 200")
	}
	if !strings.Contains(string(body), "Hello World") {
		t.Error("Body not returned as expected")
	}
	app.API["GET"][0].Scenarios[0].Response.Payload.Type = "file"
	app.API["GET"][0].Scenarios[0].Response.Payload.Data = "./framework_watcher.go"
	server.SetApp(app)
	req = httptest.NewRequest("GET", "/get", nil)
	req.Header.Set("Accept", "application/json")
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	fileData, _ := ioutil.ReadFile("./framework_watcher.go")
	if (resp.StatusCode) != 200 {
		t.Error("Status code should be 200")
	}
	if !bytes.Equal(body, fileData) {
		t.Error("File reads completed")
	}

}
