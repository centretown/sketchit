package storage

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/centretown/sketchit/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

func TestConnnect(t *testing.T) {
	var (
		client *mongo.Client
		err    error
	)

	client, err = mongo.NewClient(options.Client().ApplyURI(testURI))
	if err != nil {
		t.Fatalf("error NewClient %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("error Connect %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Fatalf("error Ping %v", err)
	}
	t.Log("OK")
	return
}

var testDevice = &api.Device{
	Label: "esp32-01",
	Model: "ESP32",
	Ip:    "192.168.1.200",
	Port:  "esp32-01",
	Pins: []*api.Device_Pin{
		{Id: 2, Label: "LED", Purpose: "activity indicator"},
		{Id: 5, Label: "TX", Purpose: "soft serial transmitter"},
		{Id: 6, Label: "RX", Purpose: "soft serial receiver"},
	},
}

func TestApiJson(t *testing.T) {
	b, err := json.MarshalIndent(testDevice, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal %v", err)
	}

	t.Logf("%s", string(b))

	deviceB := &api.Device{}
	err = json.Unmarshal(b, deviceB)
	if err != nil {
		t.Fatalf("json.Marshal %v", err)
	}

	if !compareDevices(testDevice, deviceB) {
		t.Fatal("compareDevices failed")
	}

	t.Log("compareDevices passed")
}

func compareDevices(a, b *api.Device) (cmp bool) {
	cmp = a.Label == b.Label &&
		a.Model == b.Model &&
		a.Ip == b.Ip &&
		a.Port == b.Port
	if cmp == false {
		return
	}

	cmp = len(a.Pins) == len(b.Pins)
	if cmp == false {
		return
	}

	for i, p := range a.Pins {
		q := b.Pins[i]
		cmp = p.Id == q.Id &&
			p.Label == q.Label &&
			p.Purpose == q.Purpose
		if cmp == false {
			return
		}
	}
	return
}

func TestMongoStorageProvider(t *testing.T) {
	_, err := MongoStorageProviderNew(testURI, "sketchit-test")
	if err != nil {
		t.Fatalf("MongoStorageProviderNew: %v", err)
	}
}
