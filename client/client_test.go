package main

import (
	"testing"

	"github.com/centretown/sketchit/api"
	"golang.org/x/net/context"
)

func TestCrud(t *testing.T) {
	auth := &Authentication{
		Login:    "testing",
		Password: "test",
	}

	conn, err := connect(SnakeOil, auth)
	if err != nil {
		t.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewDevicesClient(conn)
	ctx := context.Background()
	response, err := client.SayHello(ctx, &api.PingMessage{Greeting: "foo"})
	if err != nil {
		t.Fatalf("Error when calling SayHello: %s", err)
	}
	t.Logf("Response from server: %s", response.Greeting)

	req := &api.ListDevicesRequest{Parent: "domains/cottage"}
	res, err := client.List(ctx, req)
	if err != nil {
		t.Fatalf("Error when calling List: %s", err)
	}
	t.Logf("Response from server: %s\n\n", res.Devices)

	dreq := &api.GetDeviceRequest{Name: "domains/work/devices/esp32-01"}
	device, err := client.Get(ctx, dreq)
	if err != nil {
		t.Fatalf("Error when calling Get: %s", err)
	}
	t.Logf("Response from server: %s\n\n", device)

	req = &api.ListDevicesRequest{Parent: "domains/"}
	res, err = client.List(ctx, req)
	if err != nil {
		t.Fatalf("Error when calling List: %s", err)
	}
	t.Logf("Response from server: %s\n\n", res.Devices)

}
