package api

import (
	"github.com/golang/glog"
	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// ListDevices -
func (h *RequestHandler) ListDevices(ctx context.Context, req *ListDevicesRequest) (response *ListDevicesResponse, err error) {
	var devices []*Device
	glog.Infof("LIST message %s", req.Parent)
	devices, err = h.store.ListDevices(ctx, req.Parent)
	response = &ListDevicesResponse{
		Devices: devices,
	}
	return
}

// GetDevice -
func (h *RequestHandler) GetDevice(ctx context.Context, req *GetDeviceRequest) (device *Device, err error) {
	glog.Infof("GET message name=%s", req.Name)
	device, err = h.store.GetDevice(ctx, req.Name)
	return
}

// CreateDevice -
func (h *RequestHandler) CreateDevice(ctx context.Context, req *CreateDeviceRequest) (device *Device, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	device, err = h.store.CreateDevice(ctx, req.Parent, req.Device)
	return
}

// UpdateDevice -
func (h *RequestHandler) UpdateDevice(ctx context.Context, req *UpdateDeviceRequest) (device *Device, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	device, err = h.store.UpdateDevice(ctx, req.Name, req.Device)
	return
}

// DeleteDevice -
func (h *RequestHandler) DeleteDevice(ctx context.Context, req *DeleteDeviceRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.store.DeleteDevice(ctx, req.Name)
	return &empty.Empty{}, err
}
