package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/samtholiya/mockServer/common"

	"gopkg.in/yaml.v2"
)

type MockData struct {
	Name   string
	Number int
	Desc   string
}

func TestYamlSchema(t *testing.T) {
	app := App{}
	app.Post = append(app.Post, API{})
	app.Post[0].Scenarios = append(app.Post[0].Scenarios, Scenario{})
	app.Post[0].Scenarios[0].Request.Header = make(map[string][]string)
	app.Post[0].Scenarios[0].Request.Header["Key1"] = []string{"Value1"}
	app.Post[0].Scenarios[0].Request.Header["Key2"] = []string{"Value2"}

	app.Post[0].Scenarios[0].Request.Query = make(map[string][]string)
	app.Post[0].Scenarios[0].Request.Query["Key1"] = []string{"Value1"}
	app.Post[0].Scenarios[0].Request.Query["Key2"] = []string{"Value2"}
	mockData := MockData{
		Name:   "Sam",
		Number: 23,
		Desc:   "Something really nothing and then \\w+ \\d+ [] {}%^&*",
	}
	app.Post[0].Scenarios[0].Response.Payload.Type = "json"
	temp, _ := json.Marshal(mockData)
	app.Post[0].Scenarios[0].Response.Payload.Data = string(temp)
	dataFound, _ := yaml.Marshal(app)
	if os.Getenv("DEBUG") == "true" {
		file, err := os.Create("./sample_generated.yaml")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = file.WriteString(string(dataFound))
		if err != nil {
			t.Error(err)
			return
		}
		file.Close()
	}
	dataActual, err := ioutil.ReadFile(common.GetEnv("SAMPLE_YAML", "../../sample_generated.yaml"))
	if err != nil {
		t.Error(err)
		return
	}

	if string(dataActual) != string(dataFound) {
		t.Errorf("Yaml created does not match the format mentioned in sample")
	}
}

func TestJsonFeature(t *testing.T) {
	dataActual, err := ioutil.ReadFile(common.GetEnv("SAMPLE_YAML", "../../sample_generated.yaml"))
	if err != nil {
		t.Error(err)
		return
	}
	app := &App{}
	err = yaml.Unmarshal(dataActual, app)
	if err != nil {
		t.Error(err)
		return
	}
	tempData := MockData{}
	mockData := MockData{
		Name:   "Sam",
		Number: 23,
		Desc:   "Something really nothing and then \\w+ \\d+ [] {}%^&*",
	}
	_ = json.Unmarshal([]byte(app.Post[0].Scenarios[0].Response.Payload.Data), &tempData)
	if tempData != mockData {
		t.Error("Data did not match")
	}
}
