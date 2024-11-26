package knight

// Value is the interface implemented by all types within our Knight implementation.
//
// This not only includes the `Integer`, `String`, `Boolean`, `Null`, and `List` types that the spec
// defines, but also the `Variable` and the `Ast` types as well. See each type for more details.
type Value interface {
	// Run executes the value, returning the result or whatever error may have occurred.
	Run() (Value, error)

	// Dump writes a debugging representation of the value to stdout.
	Dump()

	// ToBoolean coerces to a `Boolean`.
	ToBoolean() (Boolean, error)

	// ToInteger coerces to a `Integer`.
	ToInteger() (Integer, error)

	// ToString coerces to a `String`.
	ToString() (String, error)

	// ToList coerces to a `List`.
	ToList() (List, error)
}


func runToInteger(value Value) (Integer, error) {
	ran, err := value.Run()
	if err != nil {
		return 0, err
	}

	return ran.ToInteger()
}

func runToString(value Value) (String, error) {
	ran, err := value.Run()
	if err != nil {
		return "", err
	}

	return ran.ToString()
}

func runToList(value Value) (List, error) {
	ran, err := value.Run()
	if err != nil {
		return nil, err
	}

	return ran.ToList()
}

func runToBoolean(value Value) (Boolean, error) {
	ran, err := value.Run()
	if err != nil {
		return false, err
	}

	return ran.ToBoolean()
}
