package storage

import (
	"context"
	"testing"

	"github.com/centretown/sketchit/api"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"
var newDevice = &api.Device{
	Domain: "home",
	Label:  "esp32-01",
	Model:  "ESP32",
	Ip:     "192.168.1.200:8080",
	Port:   "esp32-01",
	Pins: []*api.Device_Pin{
		{Id: 2, Label: "LED", Purpose: "activity indicator"},
		{Id: 5, Label: "TX", Purpose: "soft serial transmitter"},
		{Id: 6, Label: "RX", Purpose: "soft serial receiver"},
	},
}

func TestMongoStorageProvider(t *testing.T) {
	databaseName := "sketchit-test"
	mdp, err := MongoStorageProviderNew(testURI, databaseName)
	if err != nil {
		t.Fatalf("MongoStorageProviderNew: %v", err)
	}

	ctx := context.Background()
	// device, err := mdp.CreateDevice(ctx, "cottage", "esp32-01", newDevice)
	// if err != nil {
	// 	t.Fatalf("CreateDevice: %v", err)
	// }
	// t.Log(device)

	devices, err := mdp.ListDevices(ctx, "work")
	if err != nil {
		t.Fatalf("ListDevices: %v", err)
	}
	showDevices(t, devices)

	devices, err = mdp.ListDevices(ctx, "cottage")
	if err != nil {
		t.Fatalf("ListDevices: %v", err)
	}
	showDevices(t, devices)

}

func showDevices(t *testing.T, devices []*api.Device) {
	t.Log()
	for _, d := range devices {
		t.Logf("label: %v", d)
	}
}
