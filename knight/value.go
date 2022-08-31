package knight

// Value is the interface implemented by all types within Knight.
type Value interface {
	// Run executes the `Value`, returning a result or whatever errors occurred.
	Run() (Value, error)

	// Dump writes a debugging representation of the `Value` to stdout.
	Dump()
}

// Convertible is implemented by types that can be coerced to the four required types.
type Convertible interface {
	// ToBoolean coerces the `Convertible` to a `Boolean`.
	ToBoolean() Boolean

	// ToNumber coerces the `Convertible` to a `Number`.
	ToNumber() Number

	// ToText coerces the `Convertible` to a `Text`.
	ToText() Text

	// ToList coerces the `Convertible` to a `List`.
	ToList() List
}

type Texter interface {
	Text() Text
}
