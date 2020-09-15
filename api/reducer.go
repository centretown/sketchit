package api

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

// Sprinter commander's Print/Marshall functions use this interface
// to reduce in memory output
type Sprinter interface {
	Sprint(o interface{}, flags *FlagValues) (s string, err error)
}

// Marshaler -
type Marshaler func(o interface{}) (out []byte, err error)

// Marshalers - jump table
var Marshalers [Format_XML + 1]Marshaler

// MarshalYAML -
func MarshalYAML(o interface{}) (b []byte, err error) {
	b, err = yaml.Marshal(o)
	return
}

// MarshalJSON -
func MarshalJSON(o interface{}) (b []byte, err error) {
	b, err = json.MarshalIndent(o, "", "  ")
	return
}

// MarshalXML -
func MarshalXML(o interface{}) (b []byte, err error) {
	b, err = xml.MarshalIndent(o, "", "  ")
	return
}

func init() {
	Marshalers[Format_YAML] = MarshalYAML
	Marshalers[Format_JSON] = MarshalJSON
	Marshalers[Format_XML] = MarshalXML
}

// GetMarshaler returns a marshal function given the flags
func GetMarshaler(flags *FlagValues) (m Marshaler) {
	switch flags.Format {
	case Format_YAML,
		Format_JSON,
		Format_XML:
		m = Marshalers[flags.Format]
	default:
		m = Marshalers[Format_YAML]
	}
	return
}

// Marshal -
func Marshal(o interface{}, flags *FlagValues) (b []byte, err error) {
	b, err = GetMarshaler(flags)(o)
	return
}
