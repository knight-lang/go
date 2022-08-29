package knight

type Value interface {
	Run() (Value, error)
	Dump()
}

type Literal interface {
	Bool() bool
	Int() int
	String() string
	List() []Value
}
