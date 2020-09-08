package request

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// StorageProvider -
type StorageProvider interface {
	CreateDevice(ctx context.Context, parent string, device *api.Device) (*api.Device, error)
	GetDevice(ctx context.Context, name string) (*api.Device, error)
	ListDevices(ctx context.Context, parent string) ([]*api.Device, error)
	UpdateDevice(ctx context.Context, parent string, device *api.Device) (*api.Device, error)
	DeleteDevice(ctx context.Context, name string) error

	CreateProcess(ctx context.Context, parent string, process *api.Process) (*api.Process, error)
	GetProcess(ctx context.Context, name string) (*api.Process, error)
	ListProcesses(ctx context.Context, parent string) ([]*api.Process, error)
	UpdateProcess(ctx context.Context, parent string, device *api.Process) (*api.Process, error)
	DeleteProcess(ctx context.Context, name string) error

	ListCollections(ctx context.Context, name string) ([]*api.Collection, error)
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
func (h *RequestHandler) SayHello(ctx context.Context, in *api.PingMessage) (*api.PingMessage, error) {
	glog.Infof("SayHello message %s", in.Greeting)
	return &api.PingMessage{Greeting: "bar"}, nil
}

// ListCollections -
func (h *RequestHandler) ListCollections(ctx context.Context,
	in *api.ListCollectionsRequest) (names *api.ListCollectionsResponse, err error) {

	glog.Info("LIST COLLECTIONS")
	collections, err := h.store.ListCollections(ctx, in.Name)
	if err != nil {
		return
	}
	return &api.ListCollectionsResponse{Collections: collections}, nil
}
