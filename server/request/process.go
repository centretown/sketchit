package request

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// ListSketches -
func (h *RequestHandler) ListSketches(ctx context.Context,
	req *api.ListSketchesRequest) (response *api.ListSketchesResponse, err error) {
	var Sketches []*api.Sketch
	glog.Infof("LIST message %s", req.Parent)
	Sketches, err = h.store.ListSketches(ctx, req.Parent)
	response = &api.ListSketchesResponse{
		Sketches: Sketches,
	}
	return
}

// GetSketch -
func (h *RequestHandler) GetSketch(ctx context.Context,
	req *api.GetSketchRequest) (Sketch *api.Sketch, err error) {
	glog.Infof("GET message name=%s", req.Name)
	Sketch, err = h.store.GetSketch(ctx, req.Name)
	return
}

// CreateSketch -
func (h *RequestHandler) CreateSketch(ctx context.Context, req *api.CreateSketchRequest) (Sketch *api.Sketch, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	Sketch, err = h.store.CreateSketch(ctx, req.Parent, req.Sketch)
	return
}

// UpdateSketch -
func (h *RequestHandler) UpdateSketch(ctx context.Context, req *api.UpdateSketchRequest) (Sketch *api.Sketch, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	Sketch, err = h.store.UpdateSketch(ctx, req.Name, req.Sketch)
	return
}

// DeleteSketch -
func (h *RequestHandler) DeleteSketch(ctx context.Context, req *api.DeleteSketchRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.store.DeleteSketch(ctx, req.Name)
	return &empty.Empty{}, err
}
