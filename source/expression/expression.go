package expression

import (
	"regexp"
	"strings"

	"golib/lang"
	"golib/log"
)

type Expression struct {
	original string
}

func NewExpression(original string) Expression {
	return Expression{original: original}
}

func (self Expression) Bool(name string) string {
	split := lang.Slice(self.original, regexp.MustCompile(`(\s*\(\s*|\s*\)\s*|\s*&&\s*|\s*\|\|\s*)`), -1)
	for i, token := range split {
		switch strings.TrimSpace(token) {
		case "&&":
			split[i] = " AND "
		case "||":
			split[i] = " OR "
		case "(":
		case ")":
		case "":
			split[i] = token
		default:
			split[i] = name + " LIKE '%" + token + "%'"
		}
	}
	result := strings.Join(split, "")
	log.Debug("[%s] -> [%s]", self.original, result)
	return result
}
