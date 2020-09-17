package cmdr

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Commander script api dispatcher
type Commander struct {
	ctx         context.Context
	conn        *grpc.ClientConn
	client      api.SketchitClient
	collections []*api.Collection

	Commands   []*Command
	Aliases    map[string]*Command
	Dictionary api.Dictionary
	Flags      Flags

	index     []string
	directory []string
	suffix    string
}

// New Commander
func New(ctx context.Context,
	conn *grpc.ClientConn,
	client api.SketchitClient) (cmdr *Commander) {

	cmdr = &Commander{ctx: ctx, client: client, conn: conn, Flags: defaultFlags}
	return
}

// Print using flags
func (cmdr *Commander) Print(o interface{}, fv Presentation) string {
	fmt.Println("print flagvalues", fv)
	format := fv.Format()
	reduction := fv.Detail()
	fmt.Println("print flagvalues", format, reduction)
	return Print(o, format, reduction)
}

func (cmdr *Commander) indexList(args ...bool) (s string) {
	s = fmt.Sprintln("Help is available for the following Topics:")
	for _, key := range cmdr.index {
		c, ok := cmdr.Aliases[key]
		if ok {
			var cmd *Command
			// first and only argument
			if len(args) > 0 {
				cmd = c
			} else {
				cmd = &Command{
					Topic:   c.Topic,
					Summary: Summary{Usage: c.Summary.Usage},
				}
			}
			s += cmdr.Print(cmd, cmdr.Flags.Values())
		}
	}
	return
}

type promptLevel int

const (
	collection promptLevel = iota
	parent
	item
)

// Prompt returns the prompt string
func (cmdr *Commander) Prompt() (p string) {
	dot := "."
	dotDir := strings.Join(cmdr.directory, dot)
	p += dot + dotDir + cmdr.suffix
	return
}

// parsing information
var (
	ErrEmpty           = errors.New("no input")
	ErrCommandNotFound = errors.New("command not found")
	ErrFlagNotFound    = errors.New("flag not found")
	ErrExit            = errors.New("exit")
	ErrFlagNotValid    = errors.New("not a valid value")
)

// parseFlags removes flag tokens from input arguments
// returns flag values and remaining arguments
// errors are flagged for invalid flags or values
func (cmdr *Commander) parseFlags(input []string) (flagValues Presentation, args []string, err error) {
	flagValues = make(Presentation)
	args = make([]string, 0, len(input))
	prefix, equals, colon := "-", "=", ":"

	for _, token := range input {
		if !strings.HasPrefix(token, prefix) {
			// not a flag
			args = append(args, token)
			continue
		}

		token = strings.ToUpper(strings.TrimPrefix(token, prefix))
		split := strings.Split(token, equals)
		if len(split) < 2 {
			split = strings.Split(strings.TrimPrefix(token, prefix), colon)
			if len(split) < 2 {
				// ignore no "=" or ":"
				continue
			}
		}

		key := split[0]
		flg, ok := cmdr.Flags[key]
		if !ok {
			// invalid key
			err = info.Inform(err, ErrFlagNotFound, fmt.Sprintf("flag: %v", key))
			return
		}

		value := split[1]
		if OneOf(value, flg.Oneof...) < 0 {
			// invalid value
			err = info.Inform(err, ErrFlagNotValid,
				fmt.Sprintf("flag: %v, value: %v", key, value))
			return
		}

		flagValues[key] = value
	}

	// fill in the missing flags
	fv := cmdr.Flags.Values()
	keys := []string{"f", "d"}
	for _, key := range keys {
		if _, ok := flagValues[key]; !ok {
			flagValues[key] = fv[key]
		}
	}
	return
}

// Parse -
func (cmdr *Commander) Parse(input string) (cmd *Command, flagValues Presentation, args []string, err error) {
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		err = ErrEmpty
		return
	}

	// not empty so at least one
	args = strings.Fields(input)
	verb := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	cmd, ok := cmdr.Aliases[verb]
	if !ok {
		err = info.Inform(err, ErrCommandNotFound, verb)
		return
	}

	flagValues, args, err = cmdr.parseFlags(args)
	if err != nil {
		return
	}
	return
}

// Build commander
func (cmdr *Commander) Build() {
	request := &api.ListCollectionsRequest{}
	response, err := cmdr.client.ListCollections(cmdr.ctx, request)
	if err != nil {
		glog.Fatalf("Error when calling ListCollections: %s", err)
	}
	cmdr.collections = response.Collections
	cmdr.Commands = make([]*Command, 0)
	cmdr.Dictionary = api.DictionaryNew(response.Collections)
	cmdr.suffix = " - "

	// build index and aliases
	buildIndex := func() {
		aliasCount := 0
		cmdr.index = make([]string, len(cmdr.Commands))
		for i, c := range cmdr.Commands {
			cmdr.index[i] = c.Topic
			// plus one for the topic
			aliasCount += len(c.Aliases) + 1
		}
		sort.Sort(sort.StringSlice(cmdr.index))
		cmdr.Aliases = make(map[string]*Command, aliasCount)
		for _, c := range cmdr.Commands {
			cmdr.Aliases[c.Topic] = c
			for _, a := range c.Aliases {
				cmdr.Aliases[a] = c
			}
		}
	}
	defer buildIndex()

	cmdr.Commands = []*Command{
		goCmd,
		helloCmd,
		helpCmd,
		flagsCmd,
		listCmd,
		getCmd,
		deleteCmd,
		exitCmd,
	}

	goCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
		cmdr.directory = directory(cmdr.directory, args...)
		return
	}

	helloCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
		message := "hello " + fmt.Sprintln(args)
		response, err := cmdr.client.SayHello(cmdr.ctx, &api.PingMessage{Greeting: "hello"})
		if err != nil {
			info.Inform(err, ErrHello, message)
			return
		}
		s = cmdr.Print(response, fv)
		return
	}

	helpCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
		if len(args) < 1 {
			s = cmdr.indexList()
		} else if args[0] == "all" {
			s = cmdr.indexList(true)
		} else {
			for _, arg := range args {
				c, ok := cmdr.Aliases[arg]
				if ok {
					s = c.help(fv)
				} else {
					err = info.Inform(err, errors.New("don't know"), arg)
					return
				}
			}
		}
		return
	}

	flagsCmd.Run = func(fv Presentation, cmdArgs ...string) (s string, err error) {
		for key, value := range fv {
			fmt.Printf("key:%s value:%s\n", key, value)
		}
		return
	}

	listCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
		arg := "/"
		if len(args) > 0 {
			arg = args[0]
		}
		//cmdr.client.
		parent := fmt.Sprintf("sectors/%s", arg)
		request := &api.ListDevicesRequest{Parent: parent}
		response, err := cmdr.client.ListDevices(cmdr.ctx, request)
		if err != nil {
			err = info.Inform(err, ErrList, args)
			return
		}
		if len(response.Devices) < 1 {
			fmt.Printf("Nothing to list for %v\n", args)
			return
		}
		s = cmdr.Print(response.Devices, fv)
		return
	}

	getCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
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
		s = cmdr.Print(device, fv)
		return
	}

	deleteCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
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

		s = fmt.Sprintf("Deleted %v\n", args)
		return
	}

	exitCmd.Run = func(fv Presentation, args ...string) (s string, err error) {
		err = ErrExit
		return
	}
}
