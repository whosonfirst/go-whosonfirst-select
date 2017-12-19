package main

import (
	"errors"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-select/parser"
	"github.com/whosonfirst/go-whosonfirst-select/query"
	"github.com/whosonfirst/go-whosonfirst-select/reader"
	"github.com/whosonfirst/go-whosonfirst-select/response"
	"log"
	"os"
)

func main() {

	var source = flag.String("source", "fs", "...")
	var input = flag.String("input", "json", "...")
	var root = flag.String("fs-root", "", "...")

	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "data", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")
	var s3_creds = flag.String("s3-credentials", "", "...")

	var st = flag.String("query", "", "...")

	flag.Parse()

	p, err := parser.NewSQLParser()

	if err != nil {
		log.Fatal(err)
	}

	c, err := p.Parse(*st)

	if err != nil {
		log.Fatal(err)
	}

	var q query.Query

	switch *input {
	case "json":
		q, err = query.NewJSONQuery()
	default:
		err = errors.New("Unknown or invalid input type")
	}

	if err != nil {
		log.Fatal(err)
	}

	var r reader.Reader

	switch *source {
	case "fs":
		r, err = reader.NewFSReader(*root)
	case "s3":

		cfg := reader.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_creds,
		}

		r, err = reader.NewS3Reader(cfg)
	default:
		err = errors.New("Unknown or invalid source")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, uri := range flag.Args() {

		fh, err := r.Read(uri)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		rs, err := q.Select(fh, c)

		if err != nil {
			log.Fatal(err)
		}

		wr, err := response.NewJSONResponse()

		if err != nil {
			log.Fatal(err)
		}

		wr.WriteResults(os.Stdout, rs)
	}

}
