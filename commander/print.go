package cmdr

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
)

// Print a string with formatting and projection flags
func Print(o interface{}, flags *api.FlagValues) (s string) {
	var (
		err error
		b   []byte
	)

	sp, ok := o.(api.Sprinter)
	if ok {
		s, err = sp.Sprint(o, flags)
	} else {
		b, err = api.Marshal(o, flags)
		s = string(b)
	}

	if err != nil {
		glog.Warning(err)
		return
	}
	return
}
