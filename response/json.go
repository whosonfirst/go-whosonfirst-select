package response

import (
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-select/results"
	"io"
)

type JSONResponse struct {
	Response
}

func NewJSONResponse() (Response, error) {

	r := JSONResponse{}

	return &r, nil
}

func (r *JSONResponse) WriteResults(wr io.Writer, rs results.ResultSet) error {

	body := `{}`

	for _, r := range rs.Results() {

		rsp, err := sjson.Set(body, r.Key(), r.Value())

		if err != nil {
			return err
		}

		body = rsp
	}

	_, err := wr.Write([]byte(body))
	return err
}
