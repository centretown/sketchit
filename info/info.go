package info

import (
	"fmt"
)

// Info wraps an Operator and Sequence with the error
type Info struct {
	Text interface{}
	Kind error
	Err  error
}

var sep = "\n"

func (e *Info) Error() string {
	return fmt.Sprintf("%v\n\t '%v'\n\t %v\n", e.Text, e.Kind, e.Err)
}

// Unwrap a parsing error
func (e *Info) Unwrap() error {
	return e.Err
}

// Is this error that kind
func (e *Info) Is(kind error) bool {
	return e.Kind == kind
}

// Inform wraps existing and id the new ones' kind
func Inform(current, kind error, text interface{}) error {
	info := Info{Text: text, Kind: kind}
	wrap := kind
	if current != nil {
		wrap = fmt.Errorf("%v %w", kind, current)
	}
	info.Err = wrap
	return &info
}
