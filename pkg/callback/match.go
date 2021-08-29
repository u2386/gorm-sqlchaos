package callback

import "reflect"

type (
	MatchBy interface {
		Match(matches map[string]Match) bool
	}

	Match               func(actual interface{}) bool
	Matcher             map[string]Match
	MatchByInterface    map[string]interface{}
	MatchByReflectValue map[string]reflect.Value
)

func (matcher Matcher) Match(by MatchBy) bool {
	return by.Match(matcher)
}

func (m MatchByInterface) Match(matches map[string]Match) bool {
	for name, match := range matches {
		v, ok := m[name]
		if !ok || !match(v) {
			return false
		}
	}
	return true
}

func (vm MatchByReflectValue) Match(matches map[string]Match) bool {
	for name, match := range matches {
		v, ok := vm[name]
		if !ok {
			return false
		}

		v = reflect.Indirect(v)
		if !match(v.Interface()) {
			return false
		}
	}
	return true
}
