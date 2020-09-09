package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

// Commander - parameters
type Commander struct {
	ctx         context.Context
	conn        *grpc.ClientConn
	client      api.SketchitClient
	collections []*api.Collection
	index       []string
	directory   struct {
		prompt     string
		collection string
		parent     string
		item       string
		field      string
	}
	Commands   map[string]*Command
	Dictionary api.Dictionary
}

// CommanderNew -
func CommanderNew(ctx context.Context,
	conn *grpc.ClientConn,
	client api.SketchitClient) (cmdr *Commander) {
	cmdr = &Commander{ctx: ctx, client: client, conn: conn}
	return
}

func (cmdr *Commander) run() {
	cmdr.build()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(cmdr.prompt())
		// Read the keyboard input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Remove the newline character.
		input = strings.TrimSuffix(input, "\n")

		// Skip an empty input.
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		verb := args[0]
		args = args[1:]
		c, ok := cmdr.Commands[verb]
		if ok == false {
			err = info.Inform(err, ErrCommandNotFound, verb)
		} else {
			err = c.f(args...)
		}

		if err != nil {
			fmt.Println(err)
		}
	}

}

// Error messsages
var (
	ErrHello           = errors.New("failed to greet server")
	ErrList            = errors.New("failed to list devices")
	ErrGet             = errors.New("failed to get device")
	ErrCommandNotFound = errors.New("command not found")
	ErrNotEnoughArgs   = errors.New("not enough arguments")
)

// Summary description of command
type Summary struct {
	Usage    string   `yaml:"Usage,omitempty" json:"Usage,omitempty"`
	Syntax   string   `yaml:"Syntax,omitempty" json:"Syntax,omitempty"`
	Examples []string `yaml:"Examples,omitempty" json:"Examples,omitempty"`
}

// Command user imperative
type Command struct {
	Topic   string  `yaml:"Topic,omitempty" json:"Topic,omitempty"`
	Summary Summary `yaml:"Summary,omitempty" json:"Summary,omitempty"`
	next    string
	f       func(...string) error
}

func (cmdr *Commander) indexList() {
	fmt.Println("Help is available for the following Topics:")
	for _, key := range cmdr.index {
		c, ok := cmdr.Commands[key]
		if ok {
			brief := &Command{
				Topic:   c.Topic,
				Summary: Summary{Usage: c.Summary.Usage},
			}
			b, err := yaml.Marshal(brief)
			if err != nil {
				err = info.Inform(err, errors.New("Index not found"), "indexList")
				glog.Warning(err)
				continue
			}
			fmt.Println(string(b))
		}
	}
}

func (c *Command) help() {
	b, err := yaml.Marshal(c)
	if err != nil {
		err = info.Inform(err, errors.New("Index not found"), "indexList")
		glog.Warning(err)
	}
	fmt.Println(string(b))
}

func (cmdr *Commander) prompt() (p string) {
	cmdr.directory.prompt = "- "
	p = cmdr.directory.prompt
	return
}

func (cmdr *Commander) build() {
	request := &api.ListCollectionsRequest{}
	response, err := cmdr.client.ListCollections(cmdr.ctx, request)
	if err != nil {
		glog.Fatalf("Error when calling ListCollections: %s", err)
	}
	cmdr.collections = response.Collections
	cmdr.Commands = make(map[string]*Command)
	cmdr.Dictionary = api.DictionaryNew(response.Collections)

	organize := func() {
		for k := range cmdr.Commands {
			cmdr.index = append(cmdr.index, k)
		}
		sort.Sort(sort.StringSlice(cmdr.index))
	}

	defer organize()

	// just say hello
	cmdr.Commands["cd"] = &Command{
		Topic: "cd",
		Summary: Summary{
			Usage:    "Change current directory.",
			Syntax:   "cd <collection> <parent> <label>",
			Examples: []string{"cd devices", "cd devices work", "cd sketches ESP32"},
		},
		f: func(args ...string) (err error) {
			return
		},
	}

	// just say hello
	cmdr.Commands["hello"] = &Command{
		Topic: "hello",
		Summary: Summary{
			Usage:    "Test server by saying hello.",
			Syntax:   "hello",
			Examples: []string{"hello"},
		},
		f: func(args ...string) (err error) {
			message := "hello" + fmt.Sprintln(args)
			response, err := cmdr.client.SayHello(cmdr.ctx, &api.PingMessage{Greeting: "hello"})
			if err != nil {
				info.Inform(err, ErrHello, message)
				return
			}
			fmt.Printf("Hello: %v\n\n", response.Greeting)
			return
		},
	}

	// command line flags
	cmdr.Commands["help"] = &Command{
		Topic: "help",
		Summary: Summary{
			Usage:    "Display help for a command.",
			Syntax:   "help <Topic> <Topic>...",
			Examples: []string{"help", "help list"},
		},
		f: func(args ...string) (err error) {
			if len(args) < 1 {
				cmdr.indexList()
			} else {
				for _, arg := range args {
					c, ok := cmdr.Commands[arg]
					if ok {
						c.help()
					} else {
						fmt.Println("dont't know ", arg)
					}
				}
			}
			return
		},
	}

	// command line flags
	cmdr.Commands["flags"] = &Command{
		Topic: "flags",
		Summary: Summary{
			Usage:    "Display command line flag usage.",
			Syntax:   "flags",
			Examples: []string{"flags"},
		},
		f: func(args ...string) (err error) {
			flag.CommandLine.Usage()
			// flag.Usage()
			return
		},
	}

	// list devices
	cmdr.Commands["list"] = &Command{
		Topic: "list",
		Summary: Summary{
			Usage:    "List collection.",
			Syntax:   "list <collection> <parent>",
			Examples: []string{"list", "list devices", "list devices work"},
		},
		f: func(args ...string) (err error) {
			arg := "/"
			if len(args) > 0 {
				arg = args[0]
			}
			//cmdr.client.
			parent := fmt.Sprintf("sectors/%s", arg)
			req := &api.ListDevicesRequest{Parent: parent}
			res, err := cmdr.client.ListDevices(cmdr.ctx, req)
			if err != nil {
				err = info.Inform(err, ErrList, args)
				return
			}

			if len(res.Devices) < 1 {
				fmt.Printf("Nothing to list for %v\n", args)
				return
			}

			for _, device := range res.Devices {
				showDevice(device)
			}
			return
		},
	}

	// get a device
	cmdr.Commands["get"] = &Command{
		Topic: "get",
		Summary: Summary{
			Usage:    "Get item details.",
			Syntax:   "get <collection> <parent> <label>",
			Examples: []string{"get devices work esp32-02"},
		},
		f: func(args ...string) (err error) {
			if len(args) < 2 {
				err = info.Inform(err, ErrNotEnoughArgs, fmt.Sprint(args))
				return
			}
			name := fmt.Sprintf("sectors/%s/devices/%s", args[0], args[1])
			dreq := &api.GetDeviceRequest{Name: name}
			device, err := cmdr.client.GetDevice(cmdr.ctx, dreq)
			if err != nil {
				err = info.Inform(err, ErrGet, args)
				return
			}

			showDevice(device)
			return
		},
	}

	// delete a device
	cmdr.Commands["delete"] = &Command{
		Topic: "delete",
		Summary: Summary{
			Usage:    "Delete an item.",
			Syntax:   "delete <collection> <parent> <label>",
			Examples: []string{"delete devices work esp32-02", "delete sketch ESP32 blink"},
		},
		f: func(args ...string) (err error) {
			if len(args) < 2 {
				err = info.Inform(err, ErrNotEnoughArgs, fmt.Sprint(args))
				return
			}
			name := fmt.Sprintf("sectors/%s/devices/%s", args[0], args[1])
			dreq := &api.DeleteDeviceRequest{Name: name}

			_, err = cmdr.client.DeleteDevice(cmdr.ctx, dreq)
			if err != nil {
				err = info.Inform(err, ErrGet, args)
				return
			}

			fmt.Printf("Deleted %v\n", args)
			return
		},
	}

	// exit the app
	cmdr.Commands["exit"] = &Command{
		Topic: "exit",
		Summary: Summary{
			Usage:    "Exit this program.",
			Syntax:   "exit",
			Examples: []string{"exit"},
		},
		f: func(args ...string) (err error) {
			cmdr.conn.Close()
			os.Exit(0)
			return
		},
	}
}

func showDevice(device *api.Device) {
	fmt.Printf("Sector: %v\tLabel: %v\tModel: %v\n", device.Sector, device.Label, device.Model)
	for _, p := range device.Pins {
		fmt.Printf("\tId: %v\tLabel: %v\tPurpose: %v\n", p.Id, p.Label, p.Purpose)
	}
}
