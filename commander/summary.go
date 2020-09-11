package cmdr

import "fmt"

// Summary description of command or flag
type Summary struct {
	Usage    string   `yaml:"Usage,omitempty" json:"Usage,omitempty"`
	Syntax   string   `yaml:"Syntax,omitempty" json:"Syntax,omitempty"`
	Examples []string `yaml:"Examples,omitempty" json:"Examples,omitempty"`
}

func (sum *Summary) String() (s string) {
	var egs string
	for _, eg := range sum.Examples {
		egs += eg + " "
	}
	s = fmt.Sprintf("%s, %s %s", sum.Usage, sum.Syntax, egs)
	return
}
