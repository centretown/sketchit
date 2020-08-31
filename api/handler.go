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
	CreateDevice(parent, name string, device *Device) (*Device, error)
	// parent "domain/home-iot/devices"
	// name "esp32-01"
	DeleteDevice(parent, name string) error
	// name "domain/home-iot/devices/esp32-01"
	GetDevice(name string) (*Device, error)
	// parent "domain/home-iot/devices"
	ListDevices(parent string) ([]*Device, error)
	// parent "domain/home-iot/devices"
	UpdateDevice(parent string, device *Device) (*Device, error)
}

// StorageHandler implements grpc DevicesServer interface
type StorageHandler struct {
	StorageProvider
}

// StorageHandlerNew - creates a
func StorageHandlerNew(provider StorageProvider) *StorageHandler {
	return &StorageHandler{StorageProvider: provider}
}

// List -
func (h *StorageHandler) List(ctx context.Context, req *ListDevicesRequest) (response *ListDevicesResponse, err error) {
	var devices []*Device
	devices, err = h.ListDevices(req.Parent)
	response = &ListDevicesResponse{
		Devices: devices,
	}
	return
}

// Get -
func (h *StorageHandler) Get(ctx context.Context, req *GetDeviceRequest) (device *Device, err error) {
	device, err = h.GetDevice(req.Name)
	return
}

// Create -
func (h *StorageHandler) Create(ctx context.Context, req *CreateDeviceRequest) (device *Device, err error) {
	device, err = h.CreateDevice(req.Parent, req.Label, req.Device)
	return
}

// Update -
func (h *StorageHandler) Update(ctx context.Context, req *UpdateDeviceRequest) (device *Device, err error) {
	device, err = h.UpdateDevice(req.Parent, req.Device)
	return
}

// Delete -
func (h *StorageHandler) Delete(ctx context.Context, req *DeleteDeviceRequest) (*empty.Empty, error) {
	err := h.DeleteDevice(req.Parent, req.Name)
	return &empty.Empty{}, err
}

// SayHello -
func (h *StorageHandler) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	glog.Infof("Receive message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}
