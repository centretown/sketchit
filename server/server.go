package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/centretown/sketchit/storage"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

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

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context, s *api.StorageHandler) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")
		if clientLogin != "john" {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		if clientPassword != "doe" {
			return "", fmt.Errorf("bad password %s", clientPassword)
		}
		log.Printf("authenticated client: %s", clientLogin)
		return "42", nil
	}
	return "", fmt.Errorf("missing credentials")
}

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := info.Server.(*api.StorageHandler)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
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
	// create a listener on TCP port 7777
	// address := fmt.Sprintf("%s:%d", "dragon", 7777)
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
	// creds, err := credentials.NewServerTLSFromFile("cert/ssl-cert-snakeoil.pem",
	// 	"cert/ssl-cert-snakeoil.key")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return info.Inform(err, ErrLoadKeys, "credentials")
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor)}

	//create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// grpcServer := grpc.NewServer()
	// attach the StorageHandler service to the server
	api.RegisterDevicesServer(grpcServer, storageHandler)
	// start the server
	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return info.Inform(err, ErrServe, "grpcServer")
	}
	return nil
}

func startRESTServer(address, grpcAddress, certFile string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}
	// Setup the client gRPC options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// Register ping
	err = api.RegisterDevicesHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}
	log.Printf("starting HTTP/1.1 REST server on %s", address)
	http.ListenAndServe(address, mux)
	return nil
}

// main start a gRPC server and waits for connection
func main() {
	grpcAddress := fmt.Sprintf("%s:%d", "dragon", 7777)
	restAddress := fmt.Sprintf("%s:%d", "dragon", 7778)
	certFile := "cert/ssl-cert-snakeoil.pem"
	keyFile := "cert/ssl-cert-snakeoil.key"

	// fire the gRPC server in a goroutine
	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// fire the REST server in a goroutine
	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// infinite loop
	log.Printf("Entering infinite loop")
	select {}
}
