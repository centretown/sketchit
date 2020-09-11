package cmdr

import (
	"errors"
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

// Commander script api dispatcher
type Commander struct {
	ctx         context.Context
	conn        *grpc.ClientConn
	client      api.SketchitClient
	collections []*api.Collection

	Commands   []*Command
	Aliases    map[string]*Command
	Dictionary api.Dictionary
	Flags      *Flags

	index     []string
	directory []string
	suffix    string
}

// New Commander
func New(ctx context.Context,
	conn *grpc.ClientConn,
	client api.SketchitClient,
	flags *Flags) (cmdr *Commander) {

	cmdr = &Commander{ctx: ctx, client: client, conn: conn, Flags: flags}
	return
}

// Print using flags
func (cmdr *Commander) Print(o interface{}) string {
	return Print(cmdr.Flags.Format.Value, o)
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
			s += cmdr.Print(cmd)
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
		lenAliases := 0
		cmdr.index = make([]string, len(cmdr.Commands))
		for i, c := range cmdr.Commands {
			cmdr.index[i] = c.Topic
			lenAliases += len(c.Aliases) + 1
		}
		sort.Sort(sort.StringSlice(cmdr.index))
		cmdr.Aliases = make(map[string]*Command, lenAliases)
		for _, c := range cmdr.Commands {
			cmdr.Aliases[c.Topic] = c
			for _, a := range c.Aliases {
				cmdr.Aliases[a] = c
			}
		}
	}
	defer buildIndex()

	// changer directory
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

	goCmd.F = func(args ...string) (s string, err error) {
		cmdr.directory = directory(cmdr.directory, args...)
		return
	}

	helloCmd.F = func(args ...string) (s string, err error) {
		message := "hello" + fmt.Sprintln(args)
		response, err := cmdr.client.SayHello(cmdr.ctx, &api.PingMessage{Greeting: "hello"})
		if err != nil {
			info.Inform(err, ErrHello, message)
			return
		}
		s = cmdr.Print(response)
		return
	}

	helpCmd.F = func(args ...string) (s string, err error) {
		if len(args) < 1 {
			s = cmdr.indexList()
		} else if args[0] == "all" {
			s = cmdr.indexList(true)
		} else {
			for _, arg := range args {
				c, ok := cmdr.Aliases[arg]
				if ok {
					s = c.help(cmdr.Flags.Format.Value)
				} else {
					err = info.Inform(err, errors.New("don't know"), arg)
					return
				}
			}
		}
		return
	}

	flagsCmd.F = func(args ...string) (s string, err error) {
		a, f := extractFlags(args)
		if len(f) > 1 {
			fmt.Printf("%v %v", a, f)
		}
		// if len(args) < 1 {
		// }
		return cmdr.Print(cmdr.Flags), nil
	}

	listCmd.F = func(args ...string) (s string, err error) {
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
		s = cmdr.Print(response.Devices)
		return
	}

	getCmd.F = func(args ...string) (s string, err error) {
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
		s = cmdr.Print(device)
		return
	}

	deleteCmd.F = func(args ...string) (s string, err error) {
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

	exitCmd.F = func(args ...string) (s string, err error) {
		cmdr.conn.Close()
		os.Exit(0)
		return
	}
}
