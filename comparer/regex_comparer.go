package comparer

import (
	"regexp"

	"github.com/samtholiya/apiMocker/common"

	"github.com/sirupsen/logrus"
)

type regexComparer struct {
	log *logrus.Logger
}

//String returns a bool as true if both the string matches
func (r regexComparer) String(pattern string, str string) bool {
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
			if temp, err := regexp.MatchString(compareFrom[k], v); !temp || err != nil {
				if err != nil {
					r.log.Error(err)
				}
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (r regexComparer) MapStringArr(compareFrom map[string][]string, compareTo map[string][]string) bool {
	if len(compareFrom) != len(compareTo) {
		return false
	}
	for k := range compareFrom {
		if v, ok := compareTo[k]; ok {
			for i := range compareFrom[k] {
				if temp, err := regexp.MatchString(compareFrom[k][i], v[i]); !temp || err != nil {
					if err != nil {
						r.log.Error(err)
					}
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

//NewRegexComparer returns a regex comparer
func NewRegexComparer() Comparer {
	return &regexComparer{
		log: common.GetLogger(),
	}
}
