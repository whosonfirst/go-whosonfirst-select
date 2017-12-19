package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-readwrite/cache"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-readwrite/utils"
	"io/ioutil"
	"log"
)

func main() {

	var source = flag.String("source", "fs", "...")
	var fs_root = flag.String("fs-root", "", "...")
	var http_root = flag.String("http-root", "", "...")

	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")
	var s3_creds = flag.String("s3-credentials", "", "...")

	var cache_source = flag.String("cache", "null", "...")

	var dump = flag.Bool("dump", false, "...")

	flag.Parse()

	var args []interface{}

	switch *source {
	case "fs":
		args = []interface{}{*fs_root}
	case "http":
		args = []interface{}{*http_root}
	case "s3":
		args = []interface{}{*s3_bucket, *s3_prefix, *s3_region, *s3_creds}
	default:
		// pass
	}

	r, err := reader.NewReaderFromSource(*source, args...)

	if err != nil {
		log.Fatal(err)
	}

	c, err := cache.NewCacheFromSource(*cache_source)

	if err != nil {
		log.Fatal(err)
	}

	cr, err := reader.NewCacheReader(r, c)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		ok, err := utils.TestReader(cr, path)

		if err != nil {
			log.Fatal("TEST", err)
		}

		log.Println(path, ok)

		fh, err := cr.Read(path)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		if *dump {

			body, err := ioutil.ReadAll(fh)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(body))
		}
	}
}
