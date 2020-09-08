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
)

// Commander - parameters
type Commander struct {
	ctx         context.Context
	conn        *grpc.ClientConn
	client      api.SketchitClient
	commands    map[string]*Command
	index       []string
	collections []*api.Collection
	directory   struct {
		prompt     string
		collection string
		parent     string
		item       string
		field      string
	}
	dictonary api.Dictionary
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
		c, ok := cmdr.commands[verb]
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

// Command user imperative
type Command struct {
	topic       string
	description string
	syntax      string
	examples    []string
	next        string
	f           func(...string) error
}

func (cmdr *Commander) indexList() {
	fmt.Println("Help is available for the following topics:")
	for _, i := range cmdr.index {
		c, ok := cmdr.commands[i]
		if ok {
			fmt.Println("    ", c.topic+": ", c.description)
		}
	}
}

func (c *Command) help() {
	fmt.Println("Topic:       ", c.topic)
	fmt.Println("Description: ", c.description)
	fmt.Println("Syntax:      ", c.syntax)
	fmt.Println("Examples:")
	for _, eg := range c.examples {
		fmt.Println("             ", eg)
	}
	fmt.Println()
}

func (cmdr *Commander) prompt() (p string) {
	cmdr.directory.prompt = "> "
	p = cmdr.directory.prompt
	return
}

func (cmdr *Commander) build() {
	creq := &api.ListCollectionsRequest{}
	cres, err := cmdr.client.ListCollections(cmdr.ctx, creq)
	if err != nil {
		glog.Fatalf("Error when calling ListCollections: %s", err)
	}
	cmdr.collections = cres.Collections
	cmdr.commands = make(map[string]*Command)

	// just say hello
	cmdr.commands["hello"] = &Command{
		topic:       "hello",
		description: "Test server by saying hello.",
		syntax:      "hello",
		examples:    []string{"hello"},
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
	cmdr.commands["help"] = &Command{
		topic:       "help",
		description: "Display help for a command.",
		syntax:      "help <topic> <topic>...",
		examples:    []string{"help", "help list"},
		f: func(args ...string) (err error) {
			if len(args) < 1 {
				cmdr.indexList()
			} else {
				for _, arg := range args {
					c, ok := cmdr.commands[arg]
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
	cmdr.commands["flags"] = &Command{
		topic:       "flags",
		description: "Display command line flag usage.",
		syntax:      "flags",
		examples:    []string{"flags"},
		f: func(args ...string) (err error) {
			flag.CommandLine.Usage()
			// flag.Usage()
			return
		},
	}

	// list devices
	cmdr.commands["list"] = &Command{
		topic:       "list",
		description: "List collection.",
		syntax:      "list <collection> <parent>",
		examples:    []string{"list", "list devices", "list devices work"},
		f: func(args ...string) (err error) {
			arg := "/"
			if len(args) > 0 {
				arg = args[0]
			}
			//cmdr.client.
			parent := fmt.Sprintf("domains/%s", arg)
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
	cmdr.commands["get"] = &Command{
		topic:       "get",
		description: "Get item details.",
		syntax:      "get <collection> <parent> <label>",
		examples:    []string{"get devices work esp32-02"},
		f: func(args ...string) (err error) {
			if len(args) < 2 {
				err = info.Inform(err, ErrNotEnoughArgs, fmt.Sprint(args))
				return
			}
			name := fmt.Sprintf("domains/%s/devices/%s", args[0], args[1])
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
	cmdr.commands["delete"] = &Command{
		topic:       "delete",
		description: "Delete an item.",
		syntax:      "delete <collection> <parent> <label>",
		examples:    []string{"delete devices work esp32-02", "delete process ESP32 blink"},
		f: func(args ...string) (err error) {
			if len(args) < 2 {
				err = info.Inform(err, ErrNotEnoughArgs, fmt.Sprint(args))
				return
			}
			name := fmt.Sprintf("domains/%s/devices/%s", args[0], args[1])
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
	cmdr.commands["exit"] = &Command{
		topic:       "exit",
		description: "Exit this program.",
		syntax:      "exit",
		examples:    []string{"exit"},
		f: func(args ...string) (err error) {
			cmdr.conn.Close()
			os.Exit(0)
			return
		},
	}

	for k := range cmdr.commands {
		cmdr.index = append(cmdr.index, k)
	}

	sort.Sort(sort.StringSlice(cmdr.index))
}

func showDevice(device *api.Device) {
	fmt.Printf("Domain: %v\tLabel: %v\tModel: %v\n", device.Domain, device.Label, device.Model)
	for _, p := range device.Pins {
		fmt.Printf("\tId: %v\tLabel: %v\tPurpose: %v\n", p.Id, p.Label, p.Purpose)
	}
}
