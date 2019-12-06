package skeleton

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/samtholiya/apiMocker/common"

	"gopkg.in/yaml.v2"
)

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
	dataActual, err := ioutil.ReadFile(common.GetEnv("SAMPLE_YAML", "../sample_generated.yaml"))
	if err != nil {
		t.Error(err)
		return
	}

	if string(dataActual) != string(dataFound) {
		t.Errorf("Yaml created does not match the format mentioned in sample")
	}
}
