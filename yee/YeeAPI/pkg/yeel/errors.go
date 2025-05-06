package yeel

import "fmt"

type IllegalArgumentError struct {
	field     string
	value     string
	supported string
}

func (e IllegalArgumentError) Error() string {
	return fmt.Sprintf("%s %q %s", e.field, e.value, e.supported)
}

func newIntError(field string, value int, supported string) error {
	return IllegalArgumentError{
		field:     field,
		value:     fmt.Sprintf("%d", value),
		supported: supported,
	}
}
