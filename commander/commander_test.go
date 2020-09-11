package cmdr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/auth"
	"golang.org/x/net/context"
)

var snakeoil = "../" + auth.SnakeOil

// TestDictionary -
func TestDictionary(t *testing.T) {
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

	creq := &api.ListCollectionsRequest{}
	cres, err := client.ListCollections(ctx, creq)
	if err != nil {
		t.Fatalf("Error when calling ListCollections: %s", err)
	}

	dictionary := api.DictionaryNew(cres.Collections)

	sch := dictionary[".sketches"]
	sch.SetReducer(api.ReduceNone)
	s, err := Marshal("yaml", sch)
	if err != nil {
		t.Fatalf("marshal: %s, yaml, .sketches", err)
	}
	t.Log(s)

	sch = dictionary[".devices"]
	s, err = Marshal("yaml", sch)
	if err != nil {
		t.Fatalf("marshal: %s, yaml, .devices", err)
	}
	t.Log(s)

	sch = dictionary[".sketches"]
	sch.SetReducer(api.ReduceBrief)
	s, err = Marshal("json", sch)
	if err != nil {
		t.Fatalf("marshal: %s, json, .sketches", err)
	}
	t.Log(s)

	sch = dictionary[".devices"]
	sch.SetReducer(api.ReduceSummary)
	s, err = Marshal("xml", sch)
	if err != nil {
		t.Fatalf("marshal: %s, xml, .devices", err)
	}
	t.Log(s)
}

func TestCommands(t *testing.T) {

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
	flags := GetDefaultFlags()
	cmdr := New(ctx, conn, client, flags)
	cmdr.Build()

	testHello(t, cmdr)
	testHelp(t, cmdr)
	testFlags(t, cmdr)
}

func testHello(t *testing.T, cmdr *Commander) {
	name := "hello"
	c, ok := cmdr.Aliases[name]
	if !ok {
		t.Fatalf("%s command not found.", name)
	}
	s, err := c.F("sketchit")
	if err != nil {
		t.Fatalf("hello reported: %v", err)
	}
	t.Log(s)
}

func testHelp(t *testing.T, cmdr *Commander) {
	name := "help"
	var err error
	var errFunc = func() {
		if err != nil {
			t.Fatalf("%s reported: %v", name, err)
		}
	}
	defer errFunc()

	c, ok := cmdr.Aliases[name]
	if !ok {
		t.Fatalf("%s command not found.", name)
	}

	s, err := c.F()
	if err != nil {
		return
	}

	s, err = c.F("all")
	if err != nil {
		return
	}
	t.Log(s)
	t.Log()

	s, err = c.F("list")
	if err != nil {
		return
	}
	pre := "Topic: list"
	if !strings.HasPrefix(s, pre) {
		err = fmt.Errorf("got '%v' expected '%v'", s, pre)
		return
	}
	t.Log(s)
	t.Log()

	s, err = c.F("foo")
	if err == nil {
		err = fmt.Errorf("uncaught unknown command: %v", "foo")
		return
	}
	err = nil
	t.Log(s)
	t.Log()
}

func testFlags(t *testing.T, cmdr *Commander) {
	name := "flags"
	var err error
	var s string
	var errFunc = func() {
		if err != nil {
			err = fmt.Errorf("%s reported: %v", name, err)
			return
		}
	}
	defer errFunc()

	c, ok := cmdr.Aliases[name]
	if !ok {
		err = fmt.Errorf("command not found '%s'", name)
		return
	}

	s, err = c.F()
	if err != nil {
		return
	}
	t.Log(s)
	t.Log()

	s, err = c.F("-f=yaml", "-d=full")
	if err != nil {
		return
	}
	t.Log(s)
	t.Log()
}
