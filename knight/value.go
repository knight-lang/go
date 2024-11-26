package knight

// Value is the interface implemented by all types that are usable in Knight programs.
//
// This not only includes the `Integer`, `String`, `Boolean`, `Null`, and `List` types that the spec
// defines, but also the `Variable` and the `Ast` types as well. See each type for more details.
//
// All types must define the conversion functions, however types which don't have a defined
// conversion (such as `BLOCK`'s return values) are free to always return `error`s.
type Value interface {
	// Run executes the value, returning the result or whatever error may have occurred.
	Run() (Value, error)

	// Dump writes a debugging representation of the value to stdout.
	Dump()

	// ToBoolean coerces the type to a Boolean, or returns an error if there's a problem doing so.
	ToBoolean() (Boolean, error)

	// ToInteger coerces the type to an Integer, or returns an error if there's a problem doing so.
	ToInteger() (Integer, error)

	// ToString coerces the type to a String, or returns an error if there's a problem doing so.
	ToString() (String, error)

	// ToList coerces the type to a List, or returns an error if there's a problem doing so.
	ToList() (List, error)
}

// runToBoolean is a helper function that combines Value.Run and Value.ToBoolean.
func runToBoolean(value Value) (Boolean, error) {
	ran, err := value.Run()
	if err != nil {
		return false, err
	}

	return ran.ToBoolean()
}

// runToInteger is a helper function that combines Value.Run and Value.ToInteger.
func runToInteger(value Value) (Integer, error) {
	ran, err := value.Run()
	if err != nil {
		return 0, err
	}

	return ran.ToInteger()
}

// runToString is a helper function that combines Value.Run and Value.ToString.
func runToString(value Value) (String, error) {
	ran, err := value.Run()
	if err != nil {
		return "", err
	}

	return ran.ToString()
}

// runToList is a helper function that combines Value.Run and Value.ToList.
func runToList(value Value) (List, error) {
	ran, err := value.Run()
	if err != nil {
		return nil, err
	}

	return ran.ToList()
}
