package api

import (
	"github.com/golang/glog"
	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// StorageProvider -
type StorageProvider interface {
	// parent "domain/home-iot/devices"
	// name "esp32-01"
	CreateDevice(ctx context.Context, parent string, device *Device) (*Device, error)
	// parent "domain/home-iot/devices"
	// name "esp32-01"
	DeleteDevice(ctx context.Context, name string) error
	// name "domain/home-iot/devices/esp32-01"
	GetDevice(ctx context.Context, name string) (*Device, error)
	// parent "domain/home-iot/devices"
	ListDevices(ctx context.Context, parent string) ([]*Device, error)
	// parent "domain/home-iot/devices"
	UpdateDevice(ctx context.Context, parent string, device *Device) (*Device, error)
}

// RequestHandler implements grpc DevicesServer interface
type RequestHandler struct {
	StorageProvider
}

// RequestHandlerNew - creates a
func RequestHandlerNew(provider StorageProvider) *RequestHandler {
	return &RequestHandler{StorageProvider: provider}
}

// List -
func (h *RequestHandler) List(ctx context.Context, req *ListDevicesRequest) (response *ListDevicesResponse, err error) {
	var devices []*Device
	glog.Infof("LIST message %s", req.Parent)
	devices, err = h.ListDevices(context.Background(), req.Parent)
	response = &ListDevicesResponse{
		Devices: devices,
	}
	return
}

// Get -
func (h *RequestHandler) Get(ctx context.Context, req *GetDeviceRequest) (device *Device, err error) {
	glog.Infof("GET message name=%s", req.Name)
	device, err = h.GetDevice(ctx, req.Name)
	return
}

// Create -
func (h *RequestHandler) Create(ctx context.Context, req *CreateDeviceRequest) (device *Device, err error) {
	glog.Infof("CREATE message parent=%s", req.Parent)
	device, err = h.CreateDevice(ctx, req.Parent, req.Device)
	return
}

// Update -
func (h *RequestHandler) Update(ctx context.Context, req *UpdateDeviceRequest) (device *Device, err error) {
	glog.Infof("UPDATE message name=%s", req.Name)
	device, err = h.UpdateDevice(ctx, req.Name, req.Device)
	return
}

// Delete -
func (h *RequestHandler) Delete(ctx context.Context, req *DeleteDeviceRequest) (*empty.Empty, error) {
	glog.Infof("DELETE message %s", req.Name)
	err := h.DeleteDevice(ctx, req.Name)
	return &empty.Empty{}, err
}

// SayHello -
func (h *RequestHandler) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	glog.Infof("SayHello message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}
