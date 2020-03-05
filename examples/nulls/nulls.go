// Package nulls contains wrappers for representing dummy null types.
package nulls

import "time"

type Time struct {
	Time  time.Time
	Valid bool
}

type String struct {
	String string
	Valid  bool
}

type Int struct {
	Int   int
	Valid bool
}
