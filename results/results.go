package results

type ResultSet interface {
	AddResult(Result) error
	Results() []Result
}

type Result interface {
	Key() string
	Value() interface{}
}
