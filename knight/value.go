package knight

// Value is the interface implemented by all types within our Knight implementation.
//
// This not only includes the `Integer`, `String`, `Boolean`, `Null`, and `List` types that the spec
// defines, but also the `Variable` and the `Ast` types as well. See each type for more details.
type Value interface {
	// Run executes the value, returning the result or whatever error may have occurred.
	Run() (Value, error)

	// Dump writes a debugging representation of the vaue` to stdout.
	Dump()
}

// Convertible is implemented by types that can be coerced to the four conversion" tpyes
type Convertible interface {
	// ToBoolean coerces the to a `Boolean`.
	ToBoolean() Boolean

	// ToInteger coerces the to a `Integer`.
	ToInteger() Integer

	// ToString coerces the to a `String`.
	ToString() String

	// ToList coerces the to a `List`.
	ToList() List
}
