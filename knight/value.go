package knight

// Value is the interface implemented by all types that are usable in Knight programs.
//
// This not only includes the Integer, String, Boolean, Null, and List types that the spec
// defines, but also the Variable and the Ast types as well. See each type for more details.
//
// All types must define the conversion functions, however types which don't have a defined
// conversion (such as `BLOCK`'s return values) are free to always return `error`s.
//
// NOTE: The conversion functions here return `error`s because a handful of them _are_ fallible (eg
// converting a list of `BLOCK`s to a string). However, the Knight specs say that doing any of these
// fallible conversions is undefined behaviour. As such, we _could_ do whatever we liked in these
// cases (`panic`, use a default, return an error, etc). I chose the `error` route to make error
// messages for the end-user a bit cleaner.
type Value interface {
	// Dump writes a debugging representation of the value to stdout.
	Dump()

	// Execute executes the value, returning the result or whatever error may have occurred.
	Execute() (Value, error)

	// ToBool coerces the type to a Boolean, or returns an error if there's a problem doing so.
	ToBool() (bool, error)

	// ToInt64 coerces the type to an Integer, or returns an error if there's a problem doing so.
	ToInt64() (int64, error)

	// ToString coerces the type to a String, or returns an error if there's a problem doing so.
	ToString() (string, error)

	// ToSlice coerces the type to a List, or returns an error if there's a problem doing so.
	ToSlice() ([]Value, error)
}

//
// The following are helper functions for executing Values.
//

// executeToBool is a helper function that combines Value.Execute and Value.ToBool.
func executeToBool(value Value) (bool, error) {
	ran, err := value.Execute()
	if err != nil {
		return false, err
	}

	return ran.ToBool()
}

// executeToInt64 is a helper function that combines Value.Execute and Value.ToInt64.
func executeToInt64(value Value) (int64, error) {
	ran, err := value.Execute()
	if err != nil {
		return 0, err
	}

	return ran.ToInt64()
}

// executeToString is a helper function that combines Value.Execute and Value.ToString.
func executeToString(value Value) (string, error) {
	ran, err := value.Execute()
	if err != nil {
		return "", err
	}

	return ran.ToString()
}

// executeToSlice is a helper function that combines Value.Execute and Value.ToSlice.
func executeToSlice(value Value) (List, error) {
	ran, err := value.Execute()
	if err != nil {
		return nil, err
	}

	return ran.ToSlice()
}
