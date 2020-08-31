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
	"github.com/centretown/sketchit/storage"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var testURI = "mongodb://dave:football@localhost:27017/?authSource=sketchit-test"

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

// ErrListen -
var ErrListen = errors.New("failed to listen")

// ErrStorageConnect -
var ErrStorageConnect = errors.New("failed to connect to storage provider")

// ErrLoadKeys -
var ErrLoadKeys = errors.New("could not load TLS keys")

// ErrServe -
var ErrServe = errors.New("failed to serve")

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

	// create a server instance
	storageHandler := api.StorageHandlerNew(storageProvider)

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

	// attach the StorageHandler service to the server
	api.RegisterDevicesServer(grpcServer, storageHandler)

	glog.Infof("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return info.Inform(err, ErrServe, "grpcServer")
	}
	return nil
}

// ErrLoadCert -
var ErrLoadCert = errors.New("could not load TLS certificate")

// ErrRegisterDevices -
var ErrRegisterDevices = errors.New("could not register service")

// ErrServeRest =
var ErrServeRest = errors.New("failed to serve rest")

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
	// Register ping
	err = api.RegisterDevicesHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
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

	// fire the gRPC server in a goroutine
	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			glog.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// fire the REST server in a goroutine
	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			glog.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	glog.Infof("started OK. waiting for requests...")
	// forever
	select {}
}
