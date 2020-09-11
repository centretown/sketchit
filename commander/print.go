package cmdr

import (
	"encoding/json"
	"encoding/xml"

	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// Print a string using format
func Print(format string, o interface{}) (s string) {
	s, err := Marshal(format, o)
	if err != nil {
		glog.Warning(err)
		return
	}
	// fmt.Println(s)
	return s
}

// Marshal a struct using format to determine output
func Marshal(format string, o interface{}) (s string, err error) {
	var (
		prefix = ""
		indent = "  "
		b      []byte
	)

	switch format {
	case "xml":
		b, err = xml.MarshalIndent(o, prefix, indent)
	case "json":
		b, err = json.MarshalIndent(o, prefix, indent)
	case "yaml":
		b, err = yaml.Marshal(o)
	}

	if err != nil {
		err = info.Inform(err, ErrDecode, "marshall: "+format)
		return
	}

	s = string(b)
	return
}
