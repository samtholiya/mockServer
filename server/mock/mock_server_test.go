package mock

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/samtholiya/mockServer/comparer"
	"github.com/samtholiya/mockServer/server/model"
	"github.com/samtholiya/mockServer/types"
)

func TestServerHTTPGet(t *testing.T) {
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
							StatusCode: 200,
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

func TestServerHTTPPost(t *testing.T) {
	app := model.App{
		API: map[string][]model.API{
			"POST": []model.API{model.API{
				Endpoint:    "/post",
				Description: "Here is the description",
				Scenarios: []model.Scenario{
					model.Scenario{
						Request: model.Request{
							Header: map[string][]string{"Accept": []string{"applc"}},
							Query:  map[string][]string{"id": []string{"one"}},
							Payload: model.Payload{
								Type: "json",
								Data: `{"df":"\\w+", "fd": "{{ignore_string}}"}`,
							},
						},
						Response: model.Response{
							Payload: model.Payload{
								Type: "text",
								Data: "Hello World",
							},
							StatusCode: 201,
						},
					},

					model.Scenario{
						Request: model.Request{
							Header: map[string][]string{"Accept": []string{"application/json"}},
							Query:  map[string][]string{"id": []string{"oe"}},
							Payload: model.Payload{
								Type: "json",
								Data: `{"df":"\\w+", "fd": "{{ignore_string}}"}`,
							},
						},
						Response: model.Response{
							Payload: model.Payload{
								Type: "text",
								Data: "Hello World",
							},
							StatusCode: 202,
						},
					},
					model.Scenario{
						Request: model.Request{
							Header: map[string][]string{"Accept": []string{"application/json"}},
							Query:  map[string][]string{"id": []string{"one"}},
							Payload: model.Payload{
								Type: "json",
								Data: `{"df":"\\w+", "fd": "{{ignore_string}}"}`,
							},
						},
						Response: model.Response{
							Payload: model.Payload{
								Type: "text",
								Data: "Hello World",
							},
							StatusCode: 200,
							Delay:      1,
						},
					},
				},
			},
			},
		},
	}

	server := Server{}
	server.SetWatcher(types.TestWatcher{})
	server.SetComparer(comparer.NewRegexComparer())
	server.SetApp(app)
	req := httptest.NewRequest("POST", "/post", bytes.NewReader([]byte(`{"df":"abcd", "fd": "sdfsfs"}`)))
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("id", "one")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if (resp.StatusCode) != 200 {
		t.Errorf("Status code should be 200 got %v", resp.StatusCode)
	}
	if !strings.Contains(string(body), "Hello World") {
		t.Error("Body not returned as expected")
		t.Error(string(body))
	}
}
