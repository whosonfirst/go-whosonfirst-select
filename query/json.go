package query

import (
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-select/criteria"
	"github.com/whosonfirst/go-whosonfirst-select/results"
	"io"
	"io/ioutil"
	_ "log"
)

type JSONQuery struct {
	Query
}

func NewJSONQuery() (Query, error) {

	q := JSONQuery{}

	return &q, nil
}

func (q *JSONQuery) Select(fh io.ReadCloser, c criteria.Criteria) (results.ResultSet, error) {

	rs, err := results.NewDefaultResultSet()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	for _, r := range c.Conditions() {

		path := r.Rule()
		rsp := gjson.GetBytes(body, path)

		var k string
		var v interface{}

		k = path
		v = nil

		has_alias, alias := r.Alias()

		if has_alias {
			k = alias
		}

		if rsp.Exists() {
			v = rsp.Value()
		}

		r, err := results.NewDefaultResult(k, v)

		if err != nil {
			return nil, err
		}

		err = rs.AddResult(r)

		if err != nil {
			return nil, err
		}
	}

	return rs, nil
}
