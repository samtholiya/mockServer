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
	fd["df"] = "\\d+"
	ff, _ := json.Marshal(fd)
	result := NewRegexComparer().JSONString(string(ff), "{\"df\": \"2342\"}")
	if !result {
		t.Error("JSON Comparer should return true")
	}
	result = NewRegexComparer().JSONString(string(ff), "{\"df\": 2342}")
	if result {
		t.Error("JSON Comparer should return false")
	}
}
