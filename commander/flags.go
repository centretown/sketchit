package cmdr

// Flags define available Commander flags
type Flags struct {
	Format  FormatFlag  `yaml:"Format,omitempty" json:"Format,omitempty"`
	Detail  DetailFlag  `yaml:"Detail,omitempty" json:"Detail,omitempty"`
	Confirm ConfirmFlag `yaml:"Confirm,omitempty" json:"Confirm,omitempty"`
}

// Flag -
type Flag struct {
	Key     string   `yaml:"Key,omitempty" json:"Key,omitempty"`
	Summary Summary  `yaml:"Summary,omitempty" json:"Summary,omitempty"`
	Oneof   []string `yaml:"Oneof,omitempty" json:"Oneof,omitempty"`
}

// FormatFlag -
type FormatFlag struct {
	Flag
	Value string `yaml:"Value,omitempty" json:"Value,omitempty"`
}

// DetailFlag -
type DetailFlag struct {
	Flag
	Value string `yaml:"Value,omitempty" json:"Value,omitempty"`
}

// ConfirmFlag -
type ConfirmFlag struct {
	Flag
	Value bool `yaml:"Value,omitempty" json:"Value,omitempty"`
}

var defaultFlags = &Flags{
	Format: FormatFlag{
		Flag: Flag{
			Key: "f",
			Summary: Summary{
				Usage:    "print format",
				Syntax:   "-f=<format>",
				Examples: []string{"-f=yaml", "-f=json", "-f=xml"},
			},
			Oneof: []string{"yaml", "json", "xml"},
		},
		Value: "yaml",
	},
	Detail: DetailFlag{
		Flag: Flag{
			Key: "d",
			Summary: Summary{
				Usage:    "print detail",
				Syntax:   "-d=<detail>",
				Examples: []string{"-d=full", "-d=summary", "-d=brief"},
			},
			Oneof: []string{"full", "summary", "brief"},
		},
		Value: "full",
	},
	Confirm: ConfirmFlag{
		Flag: Flag{
			Key: "y",
			Summary: Summary{
				Usage:    "auto-confirm yes to cautions",
				Syntax:   "-y",
				Examples: []string{},
			},
		},
		Value: false,
	},
}

var flagMap map[string]interface{}

// GetDefaultFlags -
func GetDefaultFlags() (flags *Flags) {
	flags = defaultFlags
	return
}

func oneOf(in string, choices []string) bool {
	for _, choice := range choices {
		if in == choice {
			return true
		}
	}
	return false
}

func extractFlags(in []string) (out []string, flags []string) {
	out = make([]string, len(in))
	flags = make([]string, len(in))
	for _, s := range in {
		if s[0] == '-' {
			flags = append(flags, s[1:])
		} else {
			out = append(out, s)
		}
	}
	return
}
