package api

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

// Marshaler function defintion:
type Marshaler func(o interface{}) (out []byte, err error)

// marshalers is a list of avalable formatters
// strictly ordered by presentation.Format
// value yaml=0,json=1,xml=2...
var marshalers = []Marshaler{
	func(o interface{}) ([]byte, error) {
		return yaml.Marshal(o)
	},
	func(o interface{}) (b []byte, err error) {
		return json.MarshalIndent(o, "", "  ")
	},
	func(o interface{}) ([]byte, error) {
		return xml.MarshalIndent(o, "", "  ")
	},
}

// Marshal -
func Marshal(o interface{}, presentation *Presentation) (b []byte, err error) {
	m := marshalers[Format_yaml]
	if presentation.Format > 0 &&
		presentation.Format < Format(len(Format_value)) {
		m = marshalers[presentation.Format]
	}
	if len(presentation.Projection) > 0 {
		o = Project(o, presentation.Projection...)
	}
	b, err = m(o)
	return
}
