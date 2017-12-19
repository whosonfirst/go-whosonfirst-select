package parser

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-select/criteria"
	_ "log"
	"regexp"
	"strings"
)

var re_columns *regexp.Regexp
var re_aliases *regexp.Regexp

func init() {

	re_columns = regexp.MustCompile(`^SELECT (.*)\s*(?:FROM\s+.*)?$`)
	re_aliases = regexp.MustCompile(`^(.*)\s+AS\s+(.*)$`)
}

type SQLParser struct {
	Parser
}

func NewSQLParser() (Parser, error) {

	p := SQLParser{}
	return &p, nil
}

func (p *SQLParser) Parse(s string) (criteria.Criteria, error) {

	cr, err := criteria.NewColumnarCriteria()

	if err != nil {
		return nil, err
	}

	s = strings.Trim(s, " ")

	columns_m := re_columns.FindAllStringSubmatch(s, -1)

	if len(columns_m) == 0 {
		return nil, errors.New("no matches")
	}

	str_cols := columns_m[0][1]

	for _, col := range strings.Split(str_cols, ",") {

		col = strings.Trim(col, " ")
		alias := ""

		aliases_m := re_aliases.FindAllStringSubmatch(col, -1)

		if len(aliases_m) > 0 {
			col = aliases_m[0][1]
			alias = aliases_m[0][2]
		}

		// log.Println("ADD", col, alias)

		cond, err := criteria.NewColumnarConditionWithAlias(col, alias)

		if err != nil {
			return nil, err
		}

		err = cr.AddCondition(cond)

		if err != nil {
			return nil, err
		}
	}

	return cr, nil
}
