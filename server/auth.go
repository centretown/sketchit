package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// authentication errors
var (
	ErrUnknownUser        = errors.New("unknown user")
	ErrBadPassword        = errors.New("bad password")
	ErrMissingCredentials = errors.New("missing credentials")
	ErrCastServer         = errors.New("unable to cast server")
	ErrAuth               = errors.New("unable to cast server")
)

// authenticateClient check the client credentials
func authenticateClient(ctx context.Context, s *RequestHandler) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")
		if clientLogin != "testing" {
			return "", info.Inform(nil, ErrUnknownUser, clientLogin)
		}
		if clientPassword != "test" {
			return "", info.Inform(nil, ErrBadPassword, clientPassword)
		}
		glog.Infof("authenticated client: %s", clientLogin)
		return "42", nil
	}
	return "", fmt.Errorf("missing credentials")
}

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Login" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := serverInfo.Server.(*RequestHandler)
	if !ok {
		return nil, info.Inform(nil, ErrCastServer, "unaryInterceptor")
	}
	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, info.Inform(nil, ErrAuth, "authenticateClient")
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}
