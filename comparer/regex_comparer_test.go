package comparer

import "testing"

func TestRegexComparer(t *testing.T) {
	reg := NewRegexComparer()
	result := reg.CompareString("\\w+-abc", "asdfsadf-abc")
	if !result {
		t.Error("Regex string should match")
	}
	result = reg.CompareString("\\w+-abc", "asdfsadfabc")
	if result {
		t.Error("Regex string should not match")
	}
	regexMap := map[string]string{
		"Hello": "\\w+-abc",
		"There": "\\w+-\\w+",
	}

	tempMap := map[string]string{
		"Hello": "asdfsa-abc",
		"There": "asdfasdf-asdfasdf",
	}

	if result = reg.CompareMapString(regexMap, tempMap); !result {
		t.Error("Regex map should match")
	}

	
	tempMap = map[string]string{
		"Hello": "asdfsaabc",
		"There": "asdfasdfasdfasdf",
	}

	if result = reg.CompareMapString(regexMap, tempMap); result {
		t.Error("Regex map should not match")
	}

}
