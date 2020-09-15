package cmdr

import (
	"errors"
)

// Command user imperative
type Command struct {
	Topic     string                                      `yaml:"Topic,omitempty" json:"Topic,omitempty"`
	Aliases   []string                                    `yaml:"Aliases,omitempty" json:"Aliases,omitempty"`
	Summary   Summary                                     `yaml:"Summary,omitempty" json:"Summary,omitempty"`
	Arguments []string                                    `yaml:"Arguments,omitempty" json:"Arguments,omitempty"`
	Run       func(FlagValues, ...string) (string, error) `yaml:"-" json:"-" xml:"-"`
}

func (c *Command) help(fv FlagValues) string {
	return Print(c, fv.Format(), fv.Detail())
}

// func extractArgs(in []string) (args []string) {
// 	args = make([]string, 0, len(in))
// 	for _, s := range in {
// 		if !strings.HasPrefix(s, "-") {
// 			args = append(args, s)
// 		}
// 	}
// 	return
// }

// Error messsages
var (
	ErrHello         = errors.New("failed to greet server")
	ErrList          = errors.New("failed to list devices")
	ErrGet           = errors.New("failed to get device")
	ErrNotEnoughArgs = errors.New("not enough arguments")
	ErrDecode        = errors.New("failed to decode")
)

var goCmd = &Command{
	Topic:   "go",
	Aliases: []string{"cd"},
	Summary: Summary{
		Usage:    "go path",
		Syntax:   "go <path> or go . up a level",
		Examples: []string{"cd .", "cd devices", "cd devices work", "cd .sketches.ESP32"},
	},
}

var helloCmd = &Command{
	Topic:   "hello",
	Aliases: []string{"ping"},
	Summary: Summary{
		Usage:    "Test server by saying hello.",
		Syntax:   "hello -f=<format> -d=<detail>",
		Examples: []string{"hello"},
	},
}

var helpCmd = &Command{
	Topic:   "help",
	Aliases: []string{"info"},
	Summary: Summary{
		Usage:    "help listing for a command.",
		Syntax:   "help -f=<format> -d=<detail> <Topic> <Topic>...",
		Examples: []string{"help", "help list"},
	},
}

var flagsCmd = &Command{
	Topic: "flags",
	Summary: Summary{
		Usage:  "list or modify current flag defaults.",
		Syntax: "flags -f=<format> -d=<detail>",
		Examples: []string{
			"flags",
			"flags -d -f",
			"flags -f=yaml",
			"flags -f=json -d=full",
			"flags -f=xml -d=summary",
			"flags -d=brief",
		},
	},
}

var listCmd = &Command{
	Topic:   "list",
	Aliases: []string{"ls"},
	Summary: Summary{
		Usage:    "List prints the path contents.",
		Syntax:   "list <path> -f=<format> -d=<detail>",
		Examples: []string{"list", "list devices", "list devices work"},
	},
}

var getCmd = &Command{
	Topic:   "get",
	Aliases: []string{"show"},
	Summary: Summary{
		Usage:    "Get item details.",
		Syntax:   "get -f=<format> -d=<detail> <path>",
		Examples: []string{"get devices work esp32-02"},
	},
}

var deleteCmd = &Command{
	Topic:   "delete",
	Aliases: []string{"del", "erase", "rm", "remove"},
	Summary: Summary{
		Usage:    "Delete an item.",
		Syntax:   "delete <collection> <path>",
		Examples: []string{"delete devices work esp32-02", "delete sketch ESP32 blink"},
	},
}

var exitCmd = &Command{
	Topic:   "exit",
	Aliases: []string{"x", "quit"},
	Summary: Summary{
		Usage:    "Exit this program.",
		Syntax:   "exit",
		Examples: []string{"exit"},
	},
}
