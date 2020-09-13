package cmdr

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// Print a string using format
func Print(o interface{}, format api.MarshalFormat, levels ...api.Reduction) (s string) {
	s, err := Marshal(o, format, levels...)
	if err != nil {
		glog.Warning(err)
		return
	}
	return
}

// Marshal a struct using format to determine output
func Marshal(o interface{}, format api.MarshalFormat, levels ...api.Reduction) (s string, err error) {
	var (
		prefix = ""
		indent = "  "
		b      []byte
	)

	reducer, ok := o.(api.Reducer)
	if ok {
		o = reducer.Reduce(levels...)
	}

	switch format {
	case api.XML:
		b, err = xml.MarshalIndent(o, prefix, indent)
	case api.JSON:
		b, err = json.MarshalIndent(o, prefix, indent)
	case api.YAML:
		b, err = yaml.Marshal(o)
	}

	if err != nil {
		err = info.Inform(err, ErrDecode,
			fmt.Sprintf("marshall:\n\t%v\n\tformat:%v levels: %v",
				o, format.String(), levels))
		return
	}

	s = string(b)
	return
}
