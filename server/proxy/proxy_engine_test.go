package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/samtholiya/mockServer/types"
)

func TestProxyServerForm(t *testing.T) {
	prox := NewProxyServer("https://httpbin.org", types.TestDataFormatParser{})
	req := httptest.NewRequest("POST", "/post", bytes.NewBuffer([]byte("{\"Hello\":\"World\"}")))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "form")
	w := httptest.NewRecorder()

	prox.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if (resp.StatusCode) != 200 {
		t.Error("Status code should be 200")
	}
	sBody := string(body)
	if !strings.Contains(sBody, "application/json") {
		t.Error("Post calls are not sending headers")
	}
	if !strings.Contains(sBody, "{\\\"Hello\\\":\\\"World\\\"}") {
		t.Error("Post calls are not sending data")
	}
	if err := os.RemoveAll("./request_files"); err != nil {
		t.Error(err)
	}
}
