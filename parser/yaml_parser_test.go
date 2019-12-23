package parser

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	yam := NewYamlParser()
	temp := map[string]string{
		"Here": "There",
	}
	if err := yam.WriteToFile(temp, "/df/"); err == nil {
		t.Error("Error writing file should be returned")
	}
	if err := yam.WriteToFile(temp, "./yaml_generated.yaml"); err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadFile("./yaml_generated.yaml")
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(string(data), "Here: There\n") != 0 {
		t.Error("Yaml data should match")
	}
	if err := os.Remove("./yaml_generated.yaml"); err != nil {
		t.Error("Error Removing file")
	}
}

func TestReadFromFile(t *testing.T) {
	yam := NewYamlParser()
	temp := map[string]string{
		"Here": "There",
	}
	if err := yam.WriteToFile(temp, "./yaml_generated.yaml"); err != nil {
		t.Error(err)
	}
	temp1 := make(map[string]string)
	err := yam.ReadFromFile("./yaml_generated.yaml", &temp1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(temp1, temp) {
		t.Error("File read results Did not match")
	}

	if err := os.Remove("./yaml_generated.yaml"); err != nil {
		t.Error("Error Removing file")
	}
}

func TestToObject(t *testing.T) {
	expected := map[string]string{
		"Here": "There",
	}
	real := make(map[string]string)
	yam := NewYamlParser()
	if err := yam.ToObject("Here: There", &real); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, real) {
		t.Error("File read results Did not match")
	}
}
