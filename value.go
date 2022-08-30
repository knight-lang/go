package knight

type Value interface {
	Run() (Value, error)
	Dump()
}

type Literal interface {
	ToBoolean() Boolean
	ToNumber() Number
	ToText() Text
	ToList() List
}
