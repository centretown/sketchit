package main

import (
	"github.com/centretown/sketchit/api"
	"github.com/golang/glog"
	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// ListDevices -
func (h *RequestHandler) ListDevices(ctx context.Context,
	req *api.ListDevicesRequest) (response *api.ListDevicesResponse, err error) {
	var devices []*api.Device
	glog.Infof("LIST message %s", req.Parent)
	devices, err = h.store.ListDevices(ctx, req.Parent)
	response = &api.ListDevicesResponse{
		Devices: devices,
	}
	return
}

// GetDevice -
func (h *RequestHandler) GetDevice(ctx context.Context,
	req *api.GetDeviceRequest) (device *api.Device, err error) {
	glog.Infof("GET message name=%s", req.Name)
	device, err = h.store.GetDevice(ctx, req.Name)
	return
}

// CreateDevice -
func (h *RequestHandler) CreateDevice(ctx context.Context,
	req *api.CreateDeviceRequest) (device *api.Device, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	device, err = h.store.CreateDevice(ctx, req.Parent, req.Device)
	return
}

// UpdateDevice -
func (h *RequestHandler) UpdateDevice(ctx context.Context,
	req *api.UpdateDeviceRequest) (device *api.Device, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	device, err = h.store.UpdateDevice(ctx, req.Name, req.Device)
	return
}

// DeleteDevice -
func (h *RequestHandler) DeleteDevice(ctx context.Context,
	req *api.DeleteDeviceRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.store.DeleteDevice(ctx, req.Name)
	return &empty.Empty{}, err
}
