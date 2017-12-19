package query

import (
	"github.com/whosonfirst/go-whosonfirst-select/criteria"
	"github.com/whosonfirst/go-whosonfirst-select/results"
	"io"
)

type Query interface {
	Select(io.ReadCloser, criteria.Criteria) (results.ResultSet, error)
}
