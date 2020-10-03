package api

import (
	context "context"
	"testing"

	"github.com/centretown/sketchit/auth"
)

func TestResponder(t *testing.T) {
	testAuth := &auth.Authentication{
		Login:    "testing",
		Password: "test",
	}

	// connect to self cert
	conn, err := auth.Connect("../"+auth.SnakeOil, testAuth)
	if err != nil {
		t.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := NewSketchitClient(conn)
	ctx := context.Background()
	responder := NewResponder(ctx, conn, client)

	mustWork(t, responder, "goto devices.sector1")
	mustWork(t, responder, "hello world")
	mustWork(t, responder, "hello world -f=json")
	mustWork(t, responder, "hello world -f=xml")
	mustWork(t, responder, "hello world -f=yaml")
	mustWork(t, responder, "hello world -d:full,summary,brief")
	mustWork(t, responder, "list ..devices.home -d:full,summary,brief")
	mustWork(t, responder, "list ..sketches.ESP32 -d:full,summary,brief")
	mustWork(t, responder, "list ..sketches -d:full,summary,brief -f:xml")
	mustWork(t, responder, "list ..sketches -d:full,summary,brief -f:json")
	mustWork(t, responder, "list ..sketches -d:full,summary,brief")
}

func mustWork(t *testing.T, responder *Responder, input string) {
	runner, err := responder.Parse(input)
	if err != nil {
		t.Fatalf("failed to parse: [%v] because of\n\t %v", input, err)
	}

	reply, err := runner.Run()
	if err != nil {
		t.Fatalf("failed to run: [%v] because of\n\t%v", input, err)
	}
	t.Log(responder.Prompt())
	t.Log(runner.presentation)
	t.Log(reply)
}
