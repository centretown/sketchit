package auth

import (
	"errors"

	"github.com/centretown/sketchit/info"
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

// connection errors
var (
	ErrNewClient = errors.New("failed to create client with credentials")
	ErrDial      = errors.New("failed to connect to server")
)

// SnakeOil self signed cert
var SnakeOil = "../cert/snakeoil/server.pem"

// Connect -
func Connect(pem string, auth *Authentication) (conn *grpc.ClientConn, err error) {
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(pem, "")
	if err != nil {
		err = info.Inform(err, ErrNewClient, "")
		return
	}

	// connect to server
	conn, err = grpc.Dial("dragon:7777",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(auth))
	if err != nil {
		err = info.Inform(err, ErrDial, pem)
		return
	}

	return
}
