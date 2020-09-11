package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/auth"
	cmdr "github.com/centretown/sketchit/commander"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

var flags *cmdr.Flags

func init() {
	flags = cmdr.GetDefaultFlags()
	flag.StringVar(&flags.Format.Value, flags.Format.Key,
		flags.Format.Value, flags.Format.Summary.String())

	flag.StringVar(&flags.Detail.Value, flags.Detail.Key,
		flags.Detail.Value, flags.Detail.Summary.String())

	flag.BoolVar(&flags.Confirm.Value, flags.Confirm.Key,
		flags.Confirm.Value, flags.Confirm.Summary.String())
}

func main() {
	flag.Parse()

	testAuth := &auth.Authentication{
		Login:    "testing",
		Password: "test",
	}

	// connect to self cert
	conn, err := auth.Connect(auth.SnakeOil, testAuth)
	if err != nil {
		glog.Errorf("did not connect: %s", err)
		return
	}
	defer conn.Close()

	client := api.NewSketchitClient(conn)
	ctx := context.Background()
	_, err = client.SayHello(ctx, &api.PingMessage{Greeting: ""})
	if err != nil {
		glog.Errorf("did not connect: %s", err)
		return
	}

	cmdr := cmdr.New(ctx, conn, client, flags)
	run(cmdr)
}

// ErrCommandNotFound -
var ErrCommandNotFound = errors.New("command not found")

func run(commander *cmdr.Commander) {
	commander.Build()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(commander.Prompt())
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
		s := ""
		c, ok := commander.Aliases[verb]
		if ok == false {
			err = info.Inform(err, ErrCommandNotFound, verb)
		} else {
			s, err = c.F(args...)
		}

		if err != nil {
			fmt.Println(err)
		} else if len(s) > 0 {
			fmt.Println(s)
		}
	}

}
