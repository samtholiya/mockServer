package comparer

import "testing"

func TestRegexComparer(t *testing.T) {
	reg := NewRegexComparer()
	result := reg.String("\\w+-abc", "asdfsadf-abc")
	if !result {
		t.Error("Regex string should match")
	}
	result = reg.String("\\w+-abc", "asdfsadfabc")
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

	if result = reg.MapString(regexMap, tempMap); !result {
		t.Error("Regex map should match")
	}

	tempMap = map[string]string{
		"Hello": "asdfsaabc",
		"There": "asdfasdfasdfasdf",
	}

	if result = reg.MapString(regexMap, tempMap); result {
		t.Error("Regex map should not match")
	}

	regexMapArr := map[string][]string{
		"Hello": []string{"\\w+-abc"},
		"There": []string{"\\w+-\\w+"},
	}

	tempMapArr := map[string][]string{
		"Hello": []string{"asdfsa-abc"},
		"There": []string{"asdfasdf-asdfasdf"},
	}

	if result = reg.MapStringArr(regexMapArr, tempMapArr); !result {
		t.Error("Regex map should match")
	}

}
