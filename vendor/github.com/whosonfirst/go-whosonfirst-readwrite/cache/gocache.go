package cache

import (
	"bytes"
	"errors"
	gocache "github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"sync/atomic"
	"time"
)

type GoCache struct {
	Cache
	Options   *GoCacheOptions
	cache     *gocache.Cache
	hits      int64
	misses    int64
	evictions int64
	keys      int64
}

type GoCacheOptions struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func DefaultGoCacheOptions() (*GoCacheOptions, error) {

	opts := GoCacheOptions{
		DefaultExpiration: 0 * time.Second,
		CleanupInterval:   0 * time.Second,
	}

	return &opts, nil
}

func NewGoCache(opts *GoCacheOptions) (Cache, error) {

	c := gocache.New(opts.DefaultExpiration, opts.CleanupInterval)

	lc := GoCache{
		Options:   opts,
		cache:     c,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
		keys:      0,
	}

	return &lc, nil
}

func (c *GoCache) Get(key string) (io.ReadCloser, error) {

	// to do: timings that don't slow everything down the way
	// go-whosonfirst-timer does now (20170915/thisisaaronland)

	cache, ok := c.cache.Get(key)

	if !ok {
		atomic.AddInt64(&c.misses, 1)
		return nil, errors.New("CACHE MISS")
	}

	atomic.AddInt64(&c.hits, 1)

	buf := bytes.NewReader(cache.([]byte))
	return nopCloser{buf}, nil
}

func (c *GoCache) Set(key string, fh io.ReadCloser) (io.ReadCloser, error) {

	/*

	   Assume an io.Reader is hooked up to a satellite dish receiving a message (maybe a 1TB message) from an
	   alien civilization who only transmits their message once every thousand years. There's no "rewinding"
	   that.

	   https://groups.google.com/forum/#!msg/golang-nuts/BzDAg0CFqyk/t3TvH9QV0xEJ

	*/

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	c.cache.Set(key, body, gocache.DefaultExpiration)
	atomic.AddInt64(&c.keys, 1)

	r := bytes.NewReader(body)
	return nopCloser{r}, nil
}

func (c *GoCache) Size() int64 {
	return atomic.LoadInt64(&c.keys)
}

func (c *GoCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *GoCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *GoCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
