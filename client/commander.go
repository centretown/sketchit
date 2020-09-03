package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Commander - parameters
type Commander struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client api.DevicesClient
}

func (cmdr *Commander) exit() {
	cmdr.conn.Close()
	os.Exit(0)
}

// Error messsages
var (
	ErrHello           = errors.New("failed to greet server")
	ErrList            = errors.New("failed to list devices")
	ErrGet             = errors.New("failed to get device")
	ErrCommandNotFound = errors.New("command not found")
	ErrNoPath          = errors.New("path required")
	ErrNotEnoughArgs   = errors.New("not enough arguments")
)

var commands = make(map[string]func([]string) error)

func (cmdr *Commander) init() {

	commands["hello"] = func(args []string) (err error) {
		message := "hello" + fmt.Sprintln(args)
		response, err := cmdr.client.SayHello(cmdr.ctx, &api.PingMessage{Greeting: "hello"})
		if err != nil {
			info.Inform(err, ErrHello, message)
			return
		}
		fmt.Printf("Hello: %v\n\n", response.Greeting)
		return
	}

	commands["list"] = func(args []string) (err error) {
		arg := "/"
		if len(args) > 0 {
			arg = args[0]
		}
		parent := fmt.Sprintf("domains/%s", arg)
		req := &api.ListDevicesRequest{Parent: parent}
		res, err := cmdr.client.List(cmdr.ctx, req)
		if err != nil {
			err = info.Inform(err, ErrList, args)
			return
		}

		for _, device := range res.Devices {
			showDevice(device)
		}
		return
	}

	commands["get"] = func(args []string) (err error) {
		if len(args) < 2 {
			err = info.Inform(err, ErrNotEnoughArgs, fmt.Sprint(args))
			return
		}
		name := fmt.Sprintf("domains/%s/devices/%s", args[0], args[1])
		dreq := &api.GetDeviceRequest{Name: name}
		device, err := cmdr.client.Get(cmdr.ctx, dreq)
		if err != nil {
			err = info.Inform(err, ErrGet, args)
			return
		}

		showDevice(device)
		return
	}
}

func showDevice(device *api.Device) {
	fmt.Printf("Domain: %v\tLabel: %v\tModel: %v\n", device.Domain, device.Label, device.Model)
	for _, p := range device.Pins {
		fmt.Printf("\tId: %v\tLabel: %v\tPurpose: %v\n", p.Id, p.Label, p.Purpose)
	}
}

func (cmdr *Commander) exec(input string) (err error) {
	// Split the input separate the command and the arguments.
	args := strings.Fields(input)
	if len(args) < 1 {
		return info.Inform(err, ErrNotEnoughArgs, input)
	}
	verb := args[0]
	args = args[1:]

	f, ok := commands[verb]
	if ok == false {
		err = info.Inform(err, ErrCommandNotFound, verb)
	}

	err = f(args)
	return
}

var prompt = "> "

func (cmdr *Commander) run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
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

		// Handle the execution of the input.
		if err = cmdr.exec(input); err != nil {
			fmt.Println(err)
		}
	}
}
