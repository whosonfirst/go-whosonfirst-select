package utils

import (
	"github.com/whosonfirst/go-whosonfirst-select/cache"
	"github.com/whosonfirst/go-whosonfirst-select/reader"
	"io"
)

func ReadWithCache(r reader.Reader, c cache.Cache, uri string) (io.ReadCloser, error) {

	var fh io.ReadCloser
	var err error

	fh, _ = c.Get(uri)

	if fh == nil {

		fh, err = r.Read(uri)

		if err != nil {
			return nil, err
		}

		fh, err = c.Set(uri, fh)

		if err != nil {
			return nil, err
		}
	}

	return fh, nil
}
