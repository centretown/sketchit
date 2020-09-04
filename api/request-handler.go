package api

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// StorageProvider -
type StorageProvider interface {
	CreateDevice(ctx context.Context, parent string, device *Device) (*Device, error)
	GetDevice(ctx context.Context, name string) (*Device, error)
	ListDevices(ctx context.Context, parent string) ([]*Device, error)
	UpdateDevice(ctx context.Context, parent string, device *Device) (*Device, error)
	DeleteDevice(ctx context.Context, name string) error

	CreateProcess(ctx context.Context, parent string, process *Process) (*Process, error)
	GetProcess(ctx context.Context, name string) (*Process, error)
	ListProcesses(ctx context.Context, parent string) ([]*Process, error)
	UpdateProcess(ctx context.Context, parent string, device *Process) (*Process, error)
	DeleteProcess(ctx context.Context, name string) error
}

// RequestHandler implements grpc DevicesServer interface
type RequestHandler struct {
	store StorageProvider
}

// RequestHandlerNew returns a
func RequestHandlerNew(provider StorageProvider) *RequestHandler {
	return &RequestHandler{store: provider}
}

// SayHello -
func (h *RequestHandler) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	glog.Infof("SayHello message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}
