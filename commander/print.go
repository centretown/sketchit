package cmdr

import (
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
)

// Print a string with formatting and projection flags
func Print(o interface{}, flags *api.Presentation) (s string) {
	var (
		err error
		b   []byte
	)

	sp, ok := o.(api.Presenter)
	if ok {
		s = sp.Present(flags)
		if len(s) == 0 {
			err = info.Inform(err, ErrEmpty,
				fmt.Sprintf("Presentation %v"))
			return
		}
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
