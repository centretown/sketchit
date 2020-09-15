package cmdr

import (
	"flag"
	"fmt"

	"github.com/centretown/sketchit/api"
)

// Flag -
type Flag struct {
	Key     string   `yaml:"Key,omitempty" json:"Key,omitempty"`
	Value   string   `yaml:"Value,omitempty" json:"Value,omitempty"`
	Name    string   `yaml:"Name,omitempty" json:"Name,omitempty"`
	Summary Summary  `yaml:"Summary,omitempty" json:"Summary,omitempty"`
	Oneof   []string `yaml:"Oneof,omitempty" json:"Oneof,omitempty"`
}

func (flg *Flag) String() string {
	return flg.Value
}

// Flags defines the Flag map
type Flags map[string]*Flag

// Values returns a FlagValues map
func (flgs Flags) Values() (flagValues FlagValues) {
	flagValues = make(FlagValues, len(flgs))
	for k, v := range flgs {
		flagValues[k] = v.Value
	}
	return
}

// FlagValues map the Value Accessor string value
type FlagValues map[string]string

// Format -
func (fv FlagValues) Format() (marshalFormat api.MarshalFormat) {
	format := fv["f"]
	marshalFormat = api.NewMarshalFormat(format)
	if marshalFormat == api.FormatNotFound {
		marshalFormat = api.YAML
		fmt.Println("format not found", format)
	}
	return
}

// Detail -
func (fv FlagValues) Detail() (reduction api.Reduction) {
	detail := fv["d"]
	reduction = api.NewReduction(detail)
	if reduction == api.ReductionNotFound {
		reduction = api.Full
	}
	return
}

// flag validation error kinds
const (
	format int = iota
	detail
	autoconfirm
)

var defaultFlags Flags = make(Flags)
var defaultFormat *Flag
var defaultDetail *Flag
var defaultYes *Flag

func init() {

	formatFlag := &Flag{
		Key:  "f",
		Name: "format",
		Summary: Summary{
			Usage:    "print format",
			Syntax:   "-f=<format>",
			Examples: []string{"-f=yaml", "-f=json", "-f=xml"},
		},
		Oneof: []string{"yaml", "json", "xml"},
		Value: "yaml",
	}
	defaultFlags["f"] = formatFlag
	flag.StringVar(&formatFlag.Value, formatFlag.Key,
		formatFlag.Value, formatFlag.Summary.String())

	detailFlag := &Flag{
		Key:  "d",
		Name: "detail",
		Summary: Summary{
			Usage:    "print detail",
			Syntax:   "-d=<detail>",
			Examples: []string{"-d=full", "-d=summary", "-d=brief"},
		},
		Oneof: []string{"full", "summary", "brief"},
		Value: "full",
	}
	defaultFlags["d"] = detailFlag
	flag.StringVar(&detailFlag.Value, detailFlag.Key,
		detailFlag.Value, detailFlag.Summary.String())

	autoFlag := &Flag{
		Key:  "a",
		Name: "auto",
		Summary: Summary{
			Usage:    "auto reply 'y' to confirmation queries",
			Syntax:   "a=<auto>",
			Examples: []string{"-a=y", "-a=n"},
		},
		Oneof: []string{"y", "n"},
		Value: "n",
	}
	defaultFlags["y"] = autoFlag
	defaultYes = defaultFlags["y"]
	flag.StringVar(&autoFlag.Value, autoFlag.Key,
		autoFlag.Value, autoFlag.Summary.String())
}

// OneOf returns < 0 or the matching item
func OneOf(s string, choices ...string) (choice int) {
	choice = -1
	for i, item := range choices {
		if s == item {
			choice = i
			return
		}
	}
	return
}
