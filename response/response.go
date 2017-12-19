package response

import (
	"github.com/whosonfirst/go-whosonfirst-select/results"
	"io"
)

type Response interface {
	WriteResults(io.Writer, results.ResultSet) error
}
