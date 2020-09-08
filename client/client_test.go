package main

import (
	"fmt"
	"testing"

	"github.com/centretown/sketchit/api"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

func TestCrud(t *testing.T) {
	auth := &Authentication{
		Login:    "testing",
		Password: "test",
	}

	conn, err := connect("../cert/snakeoil/server.pem", auth)
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

	req := &api.ListDevicesRequest{Parent: "domains/cottage"}
	res, err := client.ListDevices(ctx, req)
	if err != nil {
		t.Fatalf("Error when calling List: %s", err)
	}
	t.Logf("Response from server: %s\n\n", res.Devices)

	creq := &api.ListCollectionsRequest{}
	cres, err := client.ListCollections(ctx, creq)
	if err != nil {
		t.Fatalf("Error when calling ListCollections: %s", err)
	}

	dict := api.DictionaryNew(cres.Collections)

	sch := dict[".processes"]
	sch.SetReducer(api.ReduceNone)
	y, _ := yaml.Marshal(sch)
	fmt.Println(string(y))

	sch = dict[".devices"]
	y, _ = yaml.Marshal(sch)
	fmt.Println(string(y))

	sch = dict[".processes"]
	sch.SetReducer(api.ReduceName)
	y, _ = yaml.Marshal(sch)
	fmt.Println(string(y))

	sch = dict[".devices"]
	y, _ = yaml.Marshal(sch)
	fmt.Println(string(y))

}
