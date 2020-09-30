package request

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
)

// StorageProvider -
type StorageProvider interface {
	ListCollections(ctx context.Context, name string) ([]*api.Collection, error)
	GetDeputy(ctx context.Context, name string) (*api.Deputy, error)
	List(ctx context.Context, parent string) ([]*anypb.Any, error)
	Create(ctx context.Context, parent string, newAny *anypb.Any) (item *anypb.Any, err error)
	Get(ctx context.Context, name string) (item *anypb.Any, err error)
	Update(ctx context.Context, name string, patch *anypb.Any) (item *anypb.Any, err error)
	Delete(ctx context.Context, name string) (err error)
}

// Handler implements grpc DevicesServer interface
type Handler struct {
	store StorageProvider
}

// HandlerNew returns a
func HandlerNew(provider StorageProvider) *Handler {
	return &Handler{store: provider}
}

// SayHello -
func (h *Handler) SayHello(ctx context.Context, in *api.PingMessage) (*api.PingMessage, error) {
	glog.Infof("SayHello message %s", in.Greeting)
	return &api.PingMessage{Greeting: "bar"}, nil
}

// ListCollections -
func (h *Handler) ListCollections(ctx context.Context,
	in *api.ListCollectionsRequest) (names *api.ListCollectionsResponse, err error) {

	glog.Info("LIST COLLECTIONS")
	collections, err := h.store.ListCollections(ctx, "")
	if err != nil {
		return
	}
	return &api.ListCollectionsResponse{Collections: collections}, nil
}

// GetDeputy -
func (h *Handler) GetDeputy(ctx context.Context,
	req *api.GetDeputyRequest) (Deputy *api.Deputy, err error) {
	glog.Infof("GET message name=%s", req.Name)
	Deputy, err = h.store.GetDeputy(ctx, req.Name)
	return
}

// List -
func (h *Handler) List(ctx context.Context,
	req *api.ListRequest) (response *api.ListResponse, err error) {
	// var devices []*api.Device
	glog.Infof("LIST message %s", req.Parent)
	items, err := h.store.List(ctx, req.Parent)
	response = &api.ListResponse{
		Items: items,
	}
	return
}

// Get -
func (h *Handler) Get(ctx context.Context, req *api.GetRequest) (item *anypb.Any, err error) {
	glog.Infof("GET message name=%s", req.Name)
	item, err = h.store.Get(ctx, req.Name)
	return
}

// Create -
func (h *Handler) Create(ctx context.Context, req *api.CreateRequest) (item *anypb.Any, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	item, err = h.store.Create(ctx, req.Parent, req.Item)
	return
}

// Update -
func (h *Handler) Update(ctx context.Context, req *api.UpdateRequest) (item *anypb.Any, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	item, err = h.store.Update(ctx, req.Name, req.Item)
	return
}

// Delete -
func (h *Handler) Delete(ctx context.Context,
	req *api.DeleteRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.store.Delete(ctx, req.Name)
	return &empty.Empty{}, err
}
