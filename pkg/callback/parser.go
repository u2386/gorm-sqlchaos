package callback

import (
	"fmt"
	"strings"
)

// TODO: AST parser, maybe

// ParseThenStatement parses then clause to appliers.
// NOTE: then clause spec be like: column1=value1, column2=value2
func ParseThenStatement(then string) (Applier, error) {
	sets := make(map[string]string)
	segs := strings.Split(strings.TrimSpace(then), ",")
	for _, seg := range segs {
		kv := strings.Split(seg, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("then statement error:%s", then)
		}
		sets[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return sets, nil
}

// ParseWhenStatement parses when clause to matchers
// NOTE: when clause spec be like: column1=value1 AND column2>=value2
func ParseWhenStatement(when string) (Matcher, error) {
	segs := strings.Split(strings.TrimSpace(when), "AND")

	matches := make(map[string]Match)
	for _, seg := range segs {
		for _, op := range operators {
			if !strings.Contains(seg, op.Sign) {
				continue
			}

			res := strings.Split(seg, op.Sign)
			if len(res) != 2 {
				return nil, fmt.Errorf("when statement error:%s", seg)
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
