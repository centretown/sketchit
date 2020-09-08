package main

import (
	"errors"
	"flag"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// for glog
	flag.Parse()

	auth := &Authentication{
		Login:    "testing",
		Password: "test",
	}

	// connect to self cert
	conn, err := connect(SnakeOil, auth)
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

	cmdr := &Commander{ctx: ctx, client: client, conn: conn}
	cmdr.run()
}

// connection errors
var (
	ErrNewClient = errors.New("failed to create client with credentials")
	ErrDial      = errors.New("failed to connect to server")
)

// SnakeOil self signed cert
var SnakeOil = "cert/snakeoil/server.pem"

func connect(pem string, auth *Authentication) (conn *grpc.ClientConn, err error) {
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(pem, "")
	if err != nil {
		info.Inform(err, ErrNewClient, "")
		return
	}

	// connect to server
	conn, err = grpc.Dial("dragon:7777",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(auth))
	if err != nil {
		info.Inform(err, ErrDial, pem)
		return
	}

	return
}
