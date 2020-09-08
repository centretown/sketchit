package storage

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test-02"
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
	databaseName := "sketchit-test-02"
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

	// names, err := mdp.ListCollections(ctx)
	// if err != nil {
	// 	t.Fatalf("ListCollectionNames: %v", err)
	// }
	// t.Log("ListCollectionNames:", names)
	// var isTrue = true
	db := mdp.client.Database(mdp.Name)
	filter := bson.D{}
	opts := &options.ListCollectionsOptions{}
	// opts := &options.ListCollectionsOptions{NameOnly: &isTrue}
	cursor, err := db.ListCollections(ctx, filter, opts)
	if err != nil {
		err = info.Inform(err, ErrCollectionNames, "ListCollections")
		return
	}
	// var res []bson.D
	// err = cursor.All(ctx, &res)
	// if err != nil {
	// 	err = info.Inform(err, ErrCollectionNames, "ListCollections")
	// 	return
	// }
	// t.Logf("%+v\n", res)

	for cursor.Next(ctx) {
		c := &MongoCollection{}
		cursor.Decode(c)
		coll := c.makeCollection()

		fmt.Printf("Name: %s, Type: %s\n", coll.Name, coll.Type)
		// sch := c.Options.Validator.JSONSchema
		level := indent(0)
		// showMongoSchema(sch, c.Name, &level)
		showSchema(coll.Schema, &level)

	}
	t.Log()
	// t.Logf("ListCollections %+v", res)
	// names = []string{deviceCollectionName, processCollectionName}
	return
}

func showDevices(t *testing.T, devices []*api.Device) {
	t.Log()
	for _, d := range devices {
		t.Logf("label: %v", d)
	}
}

type indent int

func (i *indent) String() string {
	return strings.Repeat("  ", int(*i))
}

func showSchema(sch *api.Schema, level *indent) {
	fmt.Printf("%s%s\n", level, sch.Title)
	fmt.Printf("%s Name: %v, Type: %v Description:%v\n", level, sch.Name, sch.Type, sch.Description)
	if len(sch.Required) > 1 {
		fmt.Printf("%s Required: %v\n", level, sch.Required)
	}
	if len(sch.Enum) > 0 {
		fmt.Printf("%s Enum: %v\n", level, sch.Enum)
	}
	if len(sch.OneOf) > 1 {
		fmt.Printf("%s One of:\n", level)
		*level++
		for _, s := range sch.OneOf {
			showSchema(s, level)
		}
		*level--
	}
	*level++
	if sch.Items != nil {
		showSchema(sch.Items, level)
	}
	for _, s := range sch.Properties {
		showSchema(s, level)
	}
	*level--
}

func showMongoSchema(sch MongoSchema, name string, level *indent) {
	fmt.Printf("%s%s\n", level, sch.Title)
	fmt.Printf("%s Name: %v, Type: %v Description:%v\n", level, name, sch.BsonType, sch.Description)
	if len(sch.Required) > 1 {
		fmt.Printf("%s Required: %v\n", level, sch.Required)
	}
	if len(sch.Enum) > 0 {
		fmt.Printf("%s Enum: %v\n", level, sch.Enum)
	}
	if len(sch.OneOf) > 1 {
		fmt.Printf("%s One of:\n", level)
		*level++
		for i, s := range sch.OneOf {
			showMongoSchema(s, fmt.Sprint(i), level)
		}
		*level--
	}
	*level++
	if sch.Items != nil {
		showMongoSchema(*sch.Items, "items", level)
	}
	for k, s := range sch.Properties {
		showMongoSchema(s, k, level)
	}
	*level--
}
