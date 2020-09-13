package main

import (
	"testing"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/auth"
	"golang.org/x/net/context"
)

var snakeoil = "../" + auth.SnakeOil

// TestCrud -
func TestCrud(t *testing.T) {
	a := &auth.Authentication{
		Login:    "testing",
		Password: "test",
	}

	conn, err := auth.Connect(snakeoil, a)
	if err != nil {
		t.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewSketchitClient(conn)
	ctx := context.Background()
	response, err := client.SayHello(ctx, &api.PingMessage{Greeting: "foo"})
	if err != nil {
		t.Fatalf("Error when calling SayHello: %s", err)
	}

	t.Logf("Response from server: %s", response.Greeting)

	req := &api.ListDevicesRequest{Parent: "sectors/cottage"}
	res, err := client.ListDevices(ctx, req)
	if err != nil {
		t.Fatalf("Error when calling List: %s", err)
	}
	t.Logf("Response from server: %s\n\n", res.Devices)

}

// up-arrow ^[[A down-arrow^[[B right-arrow^[[C left-arrow^[[D

// xterm sequences:
// <esc>[A     - Up          <esc>[K     -             <esc>[U     -
// <esc>[B     - Down        <esc>[L     -             <esc>[V     -
// <esc>[C     - Right       <esc>[M     -             <esc>[W     -
// <esc>[D     - Left        <esc>[N     -             <esc>[X     -
// <esc>[E     -             <esc>[O     -             <esc>[Y     -
// <esc>[F     - End         <esc>[1P    - F1          <esc>[Z     -
// <esc>[G     - Keypad 5    <esc>[1Q    - F2
// <esc>[H     - Home        <esc>[1R    - F3
// <esc>[I     -             <esc>[1S    - F4
// <esc>[J     -             <esc>[T     -
func TestKeyboard(t *testing.T) {

}
