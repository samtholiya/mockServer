package comparer

import (
	"encoding/json"
	"testing"
)

func TestRegexStringComparer(t *testing.T) {
	reg := NewRegexComparer()
	result := reg.String("\\w+-abc", "asdfsadf-abc")
	if !result {
		t.Error("Regex string should match")
	}
	result = reg.String("\\w+-abc", "asdfsadfabc")
	if result {
		t.Error("Regex string should not match")
	}
}

func TestRegexMapComparer(t *testing.T) {
	reg := NewRegexComparer()
	regexMap := map[string]string{
		"Hello": "\\w+-abc",
		"There": "\\w+-\\w+",
	}

	tempMap := map[string]string{
		"Hello": "asdfsa-abc",
		"There": "asdfasdf-asdfasdf",
	}

	if result := reg.MapString(regexMap, tempMap); !result {
		t.Error("Regex map should match")
	}

	tempMap = map[string]string{
		"Hello": "asdfsaabc",
		"There": "asdfasdfasdfasdf",
	}

	if result := reg.MapString(regexMap, tempMap); result {
		t.Error("Regex map should not match")
	}
}

func TestRegexMapStringArrComparer(t *testing.T) {
	reg := NewRegexComparer()
	regexMapArr := map[string][]string{
		"Hello": []string{"\\w+-abc"},
		"There": []string{"\\w+-\\w+"},
	}

	tempMapArr := map[string][]string{
		"Hello": []string{"asdfsa-abc"},
		"There": []string{"asdfasdf-asdfasdf"},
	}

	if result := reg.MapStringArr(regexMapArr, tempMapArr); !result {
		t.Error("Regex map should match")
	}
	tempMapArr = map[string][]string{
		"Hello": []string{"asdfsaabc"},
		"There": []string{"asdfasdfasdfasdf"},
	}

	if result := reg.MapStringArr(regexMapArr, tempMapArr); result {
		t.Error("Regex map should not match")
	}

}

func TestJSONMatcher(t *testing.T) {
	fd := make(map[string]interface{})
	fd["df"] = `\d+`
	fd["ig"] = "{{ignore_string}}"
	fd["arr"] = []string{"hello"}
	fd["nArr"] = []float64{23.5}
	ff, _ := json.Marshal(fd)
	result := NewRegexComparer().JSONString(string(ff), "{\"nArr\": [23.5], \"arr\": [\"hello\"], \"df\": \"2342\", \"d1f\": \"23142\", \"ig\": \"{{ignore_string}}\"}")
	if !result {
		t.Error("JSON Comparer should return true")
	}
	result = NewRegexComparer().JSONString(string(ff), "{\"df\": 2342}")
	if result {
		t.Error("JSON Comparer should return false")
	}
	result = NewRegexComparer().JSONString(string(ff), `{"nArr": [22.5], "arr": ["hello"], "df": "2342", "d1f": "23142", "ig": "{{ignore_string}}"}`)
	if result {
		t.Error("JSON Comparer should return false")
	}

	fd["iMap"] = map[string]interface{}{
		"hello": "world",
	}
	temp := `
	{
		"nArr": [22.5],
		"arr": ["hello"],
		"df": "2342",
		"d1f": "23142",
		"ig": "{{ignore_string}}",
		"iMap": {"hello": "world"}
	}
	`
	result = NewRegexComparer().JSONString(string(ff), temp)
	if result {
		t.Error("JSON Comparer should return false")
	}
}

func TestRandomJson(t *testing.T) {
	temp := `
	{"Data":{"Agent":{"CustomerID":"someid","AgentID":"someagentid","ActivationID":"someactivationid","ProvisioningKey":"someProvisionKey","AgentVersion":"someVersion","Platform":"somePlatform","Version":"1.3"},"Client":{"Hostname":"docker-registry.com","OS":"CentOS Linux 7 (Core)","Architecture":"x86_64","MacAddress":"00:50:56:9f:08:b2","OsVersion":"4.15.14-1.el7.elrepo.x86_64","DockerVersion":"18.03.0-ce"},"Synchronization":{"Sequence":1,"RetryCount":0},"Status":{"Provision":{"Time":"1900-01-00T00:00:00Z","OsStatus":0,"HttpStatus":0,"RetryCount":0},"ConfigDownload":{"HttpStatus":0,"RetryCount":0},"SelfPatchDownload":{"HttpStatus":0,"RetryCount":0},"Setup":{"OsStatus":0,"HttpStatus":0,"RetryCount":1}}}}
	`
	compareFrom := `
	{"Data": {"Synchronization":{"Sequence":1}}}
	`
	result := NewRegexComparer().JSONString(compareFrom, temp)
	if !result {
		t.Error("JSON Comparer should return true")
	}
	compareFrom = `
	{"Data": {"Synchronization":{"Seqence":1}}}
	`
	result = NewRegexComparer().JSONString(compareFrom, temp)
	if result {
		t.Error("JSON Comparer should return false")
	}
}

func TestAnyNumber(t *testing.T) {
	temp := `
	{"Images":[{"imageUUID":"331bc1b3-ab48-7961-a1b6-6783c0ce4ef1","imageName":"","imageSha":"23265cf7e8f5821431c1b7add8531891db26bdc7bfd497a674244d5bfd37e21a","parentUUID":"331bc1b3-ab48-7961-a1b6-6783c0ce4ef1","parentSha":"","state":null,"event":"pull","eventOccurred":"1599697408998388200","history":[{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop)  CMD [\"--config\" \"nats-server.conf\"]","CreatedSince":"7 hours ago","ID":"sha256:23265cf7e8f5821431c1b7add8531891db26bdc7bfd497a674244d5bfd37e21a","Size":"41kB"},{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop)  ENTRYPOINT [\"C:\\\\nats-server.exe\"]","CreatedSince":"7 hours ago","ID":"<missing>","Size":"41kB"},{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop)  EXPOSE 4222 6222 8222","CreatedSince":"7 hours ago","ID":"<missing>","Size":"41kB"},{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop) COPY file:bef66f144841968228eb6875fdca1fb9c094da90455a3e05090bdd09e690e7ea in C:\\nats-server.conf ","CreatedSince":"7 hours ago","ID":"<missing>","Size":"41.6kB"},{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop) COPY file:0b6475b59bf6bfe6c1030fd41ea501af74fd46ae70fd98c58683b35f8ed498ff in C:\\nats-server.exe ","CreatedSince":"7 hours ago","ID":"<missing>","Size":"10.1MB"},{"Comment":"","CreatedAt":"7 hours","CreatedBy":"cmd /S /C #(nop)  ENV NATS_DOCKERIZED=1","CreatedSince":"7 hours ago","ID":"<missing>","Size":"41kB"},{"Comment":"","CreatedAt":"7 days","CreatedBy":"Apply image 1809-amd64","CreatedSince":"7 days ago","ID":"<missing>","Size":"251MB"}]}],"eventReported":132441710217030000}
	`
	compareFrom := `
	{"Images": [], "eventReported":"ANY_NUMBER"}
	`
	result := NewRegexComparer().JSONString(compareFrom, temp)
	t.Log(result)
	if !result {
		t.Error("JSON Comparer should return true")
	}
	compareFrom = `
	{"Data": {"Synchronization":{"Seqence":1}}}
	`
	result = NewRegexComparer().JSONString(compareFrom, temp)
	if result {
		t.Error("JSON Comparer should return false")
	}
}
