package criteria

type Criteria interface {
	AddCondition(Condition) error
	Conditions() []Condition
}

type Condition interface {
	Rule() string
	Alias() (bool, string)
}
