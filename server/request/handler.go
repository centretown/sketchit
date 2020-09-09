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

	CreateSketch(ctx context.Context, parent string, Sketch *api.Sketch) (*api.Sketch, error)
	GetSketch(ctx context.Context, name string) (*api.Sketch, error)
	ListSketches(ctx context.Context, parent string) ([]*api.Sketch, error)
	UpdateSketch(ctx context.Context, parent string, device *api.Sketch) (*api.Sketch, error)
	DeleteSketch(ctx context.Context, name string) error

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
