package callback

import (
	"fmt"
	"strings"
)

// ParseSetStatement parses set clause to appliers.
// NOTE: set clause spec be like: SET column1=value1, column2=value2
func ParseSetStatement(set string) (Applier, error) {
	if !strings.HasPrefix(set, "SET") {
		return nil, fmt.Errorf("set statement error:%s", set)
	}
	sets := make(map[string]string)
	segs := strings.Split(strings.TrimPrefix(strings.TrimSpace(set), "SET"), ",")
	for _, seg := range segs {
		kv := strings.Split(seg, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("set statement error:%s", set)
		}
		sets[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return sets, nil
}

// ParseWhereStatement parses where clause to matchers
// NOTE: where clause spec be like: WHERE column1=value1 AND column2>=value2
func ParseWhereStatement(where string) (Matcher, error) {
	where = strings.TrimSpace(where)
	if !strings.HasPrefix(where, "WHERE") {
		return nil, fmt.Errorf("where statement error:%s", where)
	}
	wheres := strings.Split(strings.TrimPrefix(where, "WHERE"), "AND")

	matches := make(map[string]Match)
	for _, where := range wheres {
		for _, op := range operators {
			if !strings.Contains(where, op.Sign) {
				continue
			}

			res := strings.Split(where, op.Sign)
			if len(res) != 2 {
				return nil, fmt.Errorf("where statement error:%s", where)
			}
			field := strings.TrimSpace(res[0])
			value := strings.TrimSpace(res[1])

			if _, exist := matches[field]; !exist {
				matches[field] = op.Fn(value)
				break
			}
			matches[field] = and(matches[field], op.Fn(value))
			break
		}
	}
	return matches, nil
}
