package comparer

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/samtholiya/mockServer/types"

	"github.com/samtholiya/mockServer/common"

	"github.com/sirupsen/logrus"
)

type regexComparer struct {
	log *logrus.Logger
}

//String returns a bool as true if both the string matches
func (r regexComparer) String(pattern string, str string) bool {
	//To be moved to a different package in future
	if strings.Compare(pattern, "{{ignore_string}}") == 0 {
		return true
	}
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		r.log.Errorf("Looks like Regex is not valid %v", err)
	}
	return matched
}

//MapString returns a bool as true if both the map matches
func (r regexComparer) MapString(compareFrom map[string]string, compareTo map[string]string) bool {
	if len(compareFrom) != len(compareTo) {
		return false
	}

	for k := range compareFrom {
		if v, ok := compareTo[k]; ok {
			if temp := r.String(compareFrom[k], v); !temp {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (r regexComparer) MapStringArr(compareFrom map[string][]string, compareTo map[string][]string) bool {
	if len(compareFrom) == 0 {
		return true
	}
	for k := range compareFrom {
		if v, ok := compareTo[k]; ok {
			for i := range compareFrom[k] {
				if temp := r.String(compareFrom[k][i], v[i]); !temp {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func (r regexComparer) JSONString(compareFrom string, compareTo string) bool {
	from := make(map[string]interface{})
	to := make(map[string]interface{})
	temp := strings.ReplaceAll(compareTo, "\x00", "")

	if err := json.Unmarshal([]byte(compareFrom), &from); err != nil {
		r.log.Errorf("%v occured in compareFrom string", err)
		return false
	}
	if err := json.Unmarshal([]byte(temp), &to); err != nil {
		r.log.Tracef("%v\n--------------\n%v", compareFrom, temp)
		r.log.Errorf("%v occured in compareTo string", err)
		return false
	}
	return r.JSONMap(from, to)
}

func (r regexComparer) JSONMap(compareFrom map[string]interface{}, compareTo map[string]interface{}) bool {
	return r.rootJSON(compareFrom, compareTo)
}

func (r regexComparer) iterMap(x map[string]interface{}, compareTo map[string]interface{}) bool {
	for k, v := range x {
		switch vv := v.(type) {
		case string:
			if val, ok := compareTo[k].(string); !ok || !r.String(vv, val) {
				return false
			}
		case float64:
			if val, ok := compareTo[k].(float64); !ok || val != vv {
				return false
			}
		default:
			if !r.rootJSON(v, compareTo[k]) {
				return false
			}
		}
	}
	return true
}

func (r regexComparer) iterSlice(x []interface{}, compareTo []interface{}) bool {
	for k, v := range x {
		switch vv := v.(type) {
		case string:
			if val, ok := compareTo[k].(string); !ok || !r.String(vv, val) {
				return false
			}
		case float64:
			if val, ok := compareTo[k].(float64); !ok || val != vv {
				return false
			}
		default:
			if !r.rootJSON(v, compareTo[k]) {
				return false
			}
		}
	}
	return true
}

func (r regexComparer) rootJSON(v interface{}, compareTo interface{}) bool {
	switch vv := v.(type) {
	case map[string]interface{}:
		if _, ok := compareTo.(map[string]interface{}); !ok {
			return false
		}
		return r.iterMap(vv, compareTo.(map[string]interface{}))

	case []interface{}:
		if _, ok := compareTo.([]interface{}); !ok {
			return false
		}
		return r.iterSlice(vv, compareTo.([]interface{}))
	default:

	}
	return true
}

//NewRegexComparer returns a regex comparer
func NewRegexComparer() types.Comparer {
	return &regexComparer{
		log: common.GetLogger(),
	}
}
