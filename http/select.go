package http

import (
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-select/cache"
	"github.com/whosonfirst/go-whosonfirst-select/parser"
	"github.com/whosonfirst/go-whosonfirst-select/query"
	"github.com/whosonfirst/go-whosonfirst-select/reader"
	"github.com/whosonfirst/go-whosonfirst-select/response"
	"github.com/whosonfirst/go-whosonfirst-select/utils"
	gohttp "net/http"
)

func SelectHandler(p parser.Parser, q query.Query, r reader.Reader, c cache.Cache) (gohttp.Handler, error) {

	sanitize_opts := sanitize.DefaultOptions()

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		uri := req.URL.Path

		st_raw := req.Header.Get("WOF-Query-Statement")

		st, err := sanitize.SanitizeString(st_raw, sanitize_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		// log.Println(uri, st)

		if st == "" {
			gohttp.Error(rsp, "Missing query statement", gohttp.StatusBadRequest)
			return
		}

		cr, err := p.Parse(st)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		fh, err := utils.ReadWithCache(r, c, uri)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		defer fh.Close()

		rs, err := q.Select(fh, cr)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		wr, err := response.NewJSONResponse()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		wr.WriteResults(rsp, rs)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
