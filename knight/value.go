package knight

import (
	"fmt"
)

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
	// ToBoolean coerces to a `Boolean`.
	ToBoolean() Boolean

	// ToInteger coerces to a `Integer`.
	ToInteger() Integer

	// ToString coerces to a `String`.
	ToString() String

	// ToList coerces to a `List`.
	ToList() List
}

func runToInteger(value Value) (Integer, error) {
	ran, err := value.Run()
	if err != nil {
		return 0, err
	}

	convertible, ok := ran.(Convertible)
	if !ok {
		return 0, fmt.Errorf("cannot convert %T to an Integer", ran)
	}

	return convertible.ToInteger(), nil
}

func runToString(value Value) (String, error) {
	ran, err := value.Run()
	if err != nil {
		return "", err
	}

	convertible, ok := ran.(Convertible)
	if !ok {
		return "", fmt.Errorf("cannot convert %T to a String", ran)
	}

	return convertible.ToString(), nil
}

func runToList(value Value) (List, error) {
	ran, err := value.Run()
	if err != nil {
		return nil, err
	}

	convertible, ok := ran.(Convertible)
	if !ok {
		return nil, fmt.Errorf("cannot convert %T to a List", ran)
	}

	return convertible.ToList(), nil
}

func runToBoolean(value Value) (Boolean, error) {
	ran, err := value.Run()
	if err != nil {
		return false, err
	}

	convertible, ok := ran.(Convertible)
	if !ok {
		return false, fmt.Errorf("cannot convert %T to a Boolean", ran)
	}

	return convertible.ToBoolean(), nil
}
