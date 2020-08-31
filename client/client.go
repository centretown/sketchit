package main

import (
	"flag"
	"log"

	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Authentication holds the login/password
type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	// for glog
	flag.Parse()
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("cert/snakeoil/server.pem", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
	// Setup the login/pass
	auth := &Authentication{
		Login:    "john",
		Password: "doe",
	}

	var conn *grpc.ClientConn
	// connect to server
	conn, err = grpc.Dial("dragon:7777",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(auth))

	if err != nil {
		glog.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewDevicesClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		glog.Fatalf("Error when calling SayHello: %s", err)
	}
	glog.Infof("Response from server: %s", response.Greeting)
}
