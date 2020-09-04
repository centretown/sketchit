package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	storage "github.com/centretown/sketchit/storage/mongo"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var testURI = "mongodb://dave:football@localhost:27017/?authSource=sketchit-test"

// GRPC errors
var (
	ErrListen         = errors.New("failed to listen")
	ErrStorageConnect = errors.New("failed to connect to storage provider")
	ErrLoadKeys       = errors.New("could not load TLS keys")
	ErrServe          = errors.New("failed to serve")
)

func startGRPCServer(address, certFile, keyFile string) error {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return info.Inform(err, ErrListen, address)
	}

	// create a storage provider
	storageProvider, err := storage.MongoStorageProviderNew(testURI, "sketchit-test")
	if err != nil {
		return info.Inform(err, ErrStorageConnect, testURI)
	}

	// create a rpc handler for our api
	RequestHandler := api.RequestHandlerNew(storageProvider)

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return info.Inform(err, ErrLoadKeys, "credentials")
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// attach the RequestHandler service to the server
	api.RegisterSketchitServer(grpcServer, RequestHandler)

	glog.Infof("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return info.Inform(err, ErrServe, "grpcServer")
	}
	return nil
}

// REST errors
var (
	ErrLoadCert        = errors.New("could not load TLS certificate")
	ErrRegisterDevices = errors.New("could not register service")
	ErrServeRest       = errors.New("failed to serve rest")
)

func startRESTServer(address, grpcAddress, certFile string) error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return info.Inform(err, ErrLoadCert, "credentials")
	}
	// Setup the client gRPC options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// Register Devices
	err = api.RegisterSketchitHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return info.Inform(err, ErrRegisterDevices, "Devices")
	}

	glog.Infof("starting HTTP/1.1 REST server on %s", address)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		return info.Inform(err, ErrServeRest, address)
	}
	return nil
}

// main start a gRPC server and waits for connection
func main() {
	// for glog
	flag.Parse()

	grpcAddress := fmt.Sprintf("%s:%d", "dragon", 7777)
	restAddress := fmt.Sprintf("%s:%d", "dragon", 7778)
	certFile := "cert/snakeoil/server.pem"
	keyFile := "cert/snakeoil/server.key"

	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			glog.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			glog.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	glog.Info("OK. waiting for requests...")
	// forever
	select {}
}
