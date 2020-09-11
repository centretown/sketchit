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
