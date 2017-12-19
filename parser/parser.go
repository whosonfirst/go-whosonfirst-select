package parser

import (
	"github.com/whosonfirst/go-whosonfirst-select/criteria"
)

type Parser interface {
	Parse(string) (criteria.Criteria, error)
}
