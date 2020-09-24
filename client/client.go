package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/auth"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func init() {
}

func main() {
	// for deputy and glog
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

	var responder = api.NewResponder(ctx, conn, client)
	run(responder)
}

func run(responder *api.Responder) {
	responder.Build()
	eof := false
	reader := bufio.NewReader(os.Stdin)
	for !eof {
		fmt.Print(".")
		// Read the keyboard input.
		input, err := reader.ReadString('\n')
		if err != nil {
			// ReadString only returns err when
			// no line feed was captured
			// treat all errors as eof and issue
			// warning when EOF not returned
			if err != io.EOF {
				glog.Error(err)
			}
			eof = true
		}

		s := ""
		runner, err := responder.Parse(input)
		if err == nil {
			s, err = runner.Run()
			if err == api.ErrExit {
				return
			}
		}

		if err != nil {
			if err != api.ErrEmpty {
				fmt.Println(err)
			}
		} else if len(s) > 0 {
			fmt.Println(s)
		}
	}
}
