package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ErrUnknownUser -
var ErrUnknownUser = errors.New("unknown user")

// ErrBadPassword -
var ErrBadPassword = errors.New("bad password")

// ErrMissingCredentials -
var ErrMissingCredentials = errors.New("missing credentials")

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context, s *api.StorageHandler) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")
		if clientLogin != "john" {
			return "", info.Inform(nil, ErrUnknownUser, clientLogin)
		}
		if clientPassword != "doe" {
			return "", info.Inform(nil, ErrBadPassword, clientPassword)
		}
		glog.Infof("authenticated client: %s", clientLogin)
		return "42", nil
	}
	return "", fmt.Errorf("missing credentials")
}

// ErrCastServer -
var ErrCastServer = errors.New("unable to cast server")

// ErrAuth -
var ErrAuth = errors.New("unable to cast server")

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := serverInfo.Server.(*api.StorageHandler)
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
