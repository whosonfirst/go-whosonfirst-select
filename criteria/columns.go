package criteria

import (
	"sync"
)

type ColumnarCriteria struct {
	Criteria
	mu         *sync.RWMutex
	conditions []Condition
}

func NewColumnarCriteria() (Criteria, error) {

	mu := new(sync.RWMutex)
	conditions := make([]Condition, 0)

	rs := ColumnarCriteria{
		conditions: conditions,
		mu:         mu,
	}

	return &rs, nil
}

func (rs *ColumnarCriteria) AddCondition(r Condition) error {
	rs.mu.Lock()
	rs.conditions = append(rs.conditions, r)
	rs.mu.Unlock()
	return nil
}

func (rs *ColumnarCriteria) Conditions() []Condition {
	rs.mu.RLock()
	conditions := rs.conditions
	rs.mu.RUnlock()
	return conditions
}

type ColumnarCondition struct {
	Condition
	column string
	alias  string
}

func NewColumnarCondition(column string) (Condition, error) {

	return NewColumnarConditionWithAlias(column, "")
}

func NewColumnarConditionWithAlias(column string, alias string) (Condition, error) {

	c := ColumnarCondition{
		column: column,
		alias:  alias,
	}

	return &c, nil
}

func (c *ColumnarCondition) Rule() string {
	return c.column
}

func (c *ColumnarCondition) Alias() (bool, string) {

	has_alias := c.alias != ""
	return has_alias, c.alias
}
