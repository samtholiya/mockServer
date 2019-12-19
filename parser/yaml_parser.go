package parser

import (
	"io/ioutil"
	"os"

	"github.com/samtholiya/mockServer/common"
	"github.com/samtholiya/mockServer/types"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type yamlParser struct {
	log *logrus.Logger
}

//WriteToFile Writes object to the file in the format
func (y yamlParser) WriteToFile(obj interface{}, filePath string) error {
	dataFound, err := yaml.Marshal(obj)
	if err != nil {
		y.log.Error(err)
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		y.log.Error(err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			y.log.Error(err)
		}
	}()
	_, err = file.WriteString(string(dataFound))
	if err != nil {
		y.log.Error(err)
		return err
	}
	return nil
}

//ReadFromFile Reads a given file to the object of the format
func (y yamlParser) ReadFromFile(filePath string, obj interface{}) error {
	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		y.log.Errorf("Error reading file %v", filePath)
		return err
	}
	err = yaml.Unmarshal(dataBytes, obj)
	if err != nil {
		y.log.Error(err)
		return err
	}
	y.log.Debugf("New model structure %v", obj)
	y.log.Infof("Model loaded from file %v", filePath)
	return nil
}

//ToObject Converts string to the object
func (y yamlParser) ToObject(fromString string, obj interface{}) error {
	return yaml.Unmarshal([]byte(fromString), obj)
}

//NewYamlParser Returns a yaml based data format parser
func NewYamlParser() types.DataFormatParser {
	yam := yamlParser{
		log: common.GetLogger(),
	}
	return yam
}
