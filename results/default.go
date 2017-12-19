package results

import (
	"sync"
)

type DefaultResultSet struct {
	ResultSet
	mu *sync.RWMutex
	r  []Result
}

func NewDefaultResultSet() (ResultSet, error) {

	mu := new(sync.RWMutex)
	r := make([]Result, 0)

	rs := DefaultResultSet{
		mu: mu,
		r:  r,
	}

	return &rs, nil
}

func (rs *DefaultResultSet) AddResult(r Result) error {

	rs.mu.Lock()
	rs.r = append(rs.r, r)
	rs.mu.Unlock()

	return nil
}

func (rs *DefaultResultSet) Results() []Result {
	rs.mu.RLock()
	results := rs.r
	rs.mu.RUnlock()
	return results
}

type DefaultResult struct {
	Result
	k string
	v interface{}
}

func NewDefaultResult(k string, v interface{}) (Result, error) {

	r := DefaultResult{
		k: k,
		v: v,
	}

	return &r, nil
}

func (r *DefaultResult) Key() string {
	return r.k
}

func (r *DefaultResult) Value() interface{} {
	return r.v
}
