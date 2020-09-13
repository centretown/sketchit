package cmdr

import (
	"encoding"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Accessor interface
type Accessor interface {
	String() string
	Addr() interface{}
	yaml.Marshaler
	yaml.Unmarshaler
	json.Marshaler
	json.Unmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

// StringAccessor -
type StringAccessor struct {
	value string
}

// String - UnmarshalText
func (sa *StringAccessor) String() string {
	return string(sa.value)
}

// Addr - return a pointer to the actual value
func (sa *StringAccessor) Addr() interface{} {
	return &sa.value
}

// MarshalText -
func (sa *StringAccessor) MarshalText() (text []byte, err error) {
	text = []byte(sa.String())
	return
}

// UnmarshalText -
func (sa *StringAccessor) UnmarshalText(text []byte) (err error) {
	sa.value = string(text)
	return
}

// MarshalYAML -
func (sa *StringAccessor) MarshalYAML() (text interface{}, err error) {
	text = sa.String()
	return
}

// UnmarshalYAML -
func (sa *StringAccessor) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&sa.value)
	if err != nil {
		return
	}
	// UnmarshalText validates
	err = sa.UnmarshalText([]byte(sa.value))
	return
}

// MarshalJSON -
func (sa *StringAccessor) MarshalJSON() (text []byte, err error) {
	text = []byte(sa.String())
	return
}

// UnmarshalJSON -
func (sa *StringAccessor) UnmarshalJSON(b []byte) error {
	// UnmarshalText validates
	return sa.UnmarshalText(b)
}

// BoolAccessor -
type BoolAccessor struct {
	value bool
}

// String -
func (ba *BoolAccessor) String() (s string) {
	s = "n"
	if ba.value {
		s = "y"
	}
	return
}

// Addr - return a pointer to the actual value
func (ba *BoolAccessor) Addr() interface{} {
	return &ba.value
}

// MarshalText -
func (ba *BoolAccessor) MarshalText() (text []byte, err error) {
	text = []byte(ba.String())
	return
}

// UnmarshalText -
func (ba *BoolAccessor) UnmarshalText(text []byte) (err error) {
	ba.value = false
	s := string(text)
	if s == "y" {
		ba.value = true
	}
	return
}

// MarshalYAML -
func (ba *BoolAccessor) MarshalYAML() (text interface{}, err error) {
	text = ba.String()
	return
}

// UnmarshalYAML -
func (ba *BoolAccessor) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(ba.value)
	return
}

// MarshalJSON -
func (ba *BoolAccessor) MarshalJSON() (text []byte, err error) {
	text = []byte(ba.String())
	return
}

// UnmarshalJSON -
func (ba *BoolAccessor) UnmarshalJSON(b []byte) error {
	return ba.UnmarshalText(b)
}
