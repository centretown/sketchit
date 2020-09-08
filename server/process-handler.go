package main

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// ListProcesses -
func (h *RequestHandler) ListProcesses(ctx context.Context,
	req *api.ListProcessesRequest) (response *api.ListProcessesResponse, err error) {
	var Processes []*api.Process
	glog.Infof("LIST message %s", req.Parent)
	Processes, err = h.store.ListProcesses(ctx, req.Parent)
	response = &api.ListProcessesResponse{
		Processes: Processes,
	}
	return
}

// GetProcess -
func (h *RequestHandler) GetProcess(ctx context.Context,
	req *api.GetProcessRequest) (Process *api.Process, err error) {
	glog.Infof("GET message name=%s", req.Name)
	Process, err = h.store.GetProcess(ctx, req.Name)
	return
}

// CreateProcess -
func (h *RequestHandler) CreateProcess(ctx context.Context, req *api.CreateProcessRequest) (Process *api.Process, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	Process, err = h.store.CreateProcess(ctx, req.Parent, req.Process)
	return
}

// UpdateProcess -
func (h *RequestHandler) UpdateProcess(ctx context.Context, req *api.UpdateProcessRequest) (Process *api.Process, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	Process, err = h.store.UpdateProcess(ctx, req.Name, req.Process)
	return
}

// DeleteProcess -
func (h *RequestHandler) DeleteProcess(ctx context.Context, req *api.DeleteProcessRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.store.DeleteProcess(ctx, req.Name)
	return &empty.Empty{}, err
}
