package types

//DataFormatParser parses data to yaml and other obj
type DataFormatParser interface {

	//WriteToFile Writes object to the file in the format
	WriteToFile(obj interface{}, filePath string) error

	//ReadFromFile Reads a given file to the object of the format
	ReadFromFile(filePath string, obj interface{}) error

	//ToObject Converts string to the object
	ToObject(fromString string, toObj interface{}) error
}

//TestDataFormatParser is empty implementation of DataFormatParser for testing purpose
type TestDataFormatParser struct{}

//WriteToFile Writes object to the file in the format
func (t TestDataFormatParser) WriteToFile(obj interface{}, filePath string) error {
	return nil
}

//ReadFromFile Reads a given file to the object of the format
func (t TestDataFormatParser) ReadFromFile(filePath string, obj interface{}) error {
	return nil
}

//ToObject Converts string to the object
func (t TestDataFormatParser) ToObject(fromString string, toObj interface{}) error {
	return nil
}
