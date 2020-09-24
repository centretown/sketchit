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
	Sector: "home",
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
var (
	databaseName = "sketchit-test-03"
	authSource   = "sketchit-test-02"
)

func TestMongoGetDeputy(t *testing.T) {
	databaseName := "sketchit-test-03"
	mdp, err := MongoStorageProviderNew(testURI, databaseName, authSource)
	if err != nil {
		t.Fatalf("MongoStorageProviderNew: %v", err)
	}

	ctx := context.Background()
	var deputy *api.Deputy
	deputy, err = mdp.GetDeputy(ctx, "Andy")
	if err != nil {
		t.Fatalf("GetDeputy: %v", err)
	}
	t.Log(showDeputy(deputy))
}

func showDeputy(deputy *api.Deputy) (s string) {
	s += fmt.Sprintln(deputy.Label, deputy.Version)
	for _, f := range deputy.Skills {
		s += fmt.Sprintln(f.Task, f.Alternates)
	}
	for _, f := range deputy.Features {
		s += fmt.Sprintln(f.Flag, f.Label)
	}
	return
}

func TestMongoListCollection(t *testing.T) {
	mdp, err := MongoStorageProviderNew(testURI, databaseName, authSource)
	if err != nil {
		t.Fatalf("MongoStorageProviderNew: %v", err)
	}

	ctx := context.Background()
	db := mdp.client.Database(mdp.Name)
	filter := bson.D{}
	opts := &options.ListCollectionsOptions{}

	cursor, err := db.ListCollections(ctx, filter, opts)
	if err != nil {
		err = info.Inform(err, ErrCollectionNames, "ListCollections")
		t.Fatalf("MongoStorageProviderNew: %v", err)
		return
	}

	t.Log("ListCollections", cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		c := &MongoCollection{}
		cursor.Decode(c)
		coll := c.MongoCollectionNew()

		fmt.Printf("Name: %s, Type: %s\n", coll.Name, coll.Type)
		level := indent(0)
		showSchema(coll.Model, &level)

	}
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

func showSchema(sch *api.Model, level *indent) {
	fmt.Printf("%s%s\n", level, sch.Title)
	fmt.Printf("%s Name: %v, Type: %v Description:%v\n", level, sch.Label, sch.Type, sch.Description)
	if len(sch.Required) > 1 {
		fmt.Printf("%s Required: %v\n", level, sch.Required)
	}
	if len(sch.Options) > 0 {
		fmt.Printf("%s Options: %v\n", level, sch.Options)
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

// devices, err := mdp.ListDevices(ctx, "work")
// if err != nil {
// 	t.Fatalf("ListDevices: %v", err)
// }
// showDevices(t, devices)

// devices, err = mdp.ListDevices(ctx, "cottage")
// if err != nil {
// 	t.Fatalf("ListDevices: %v", err)
// }
// showDevices(t, devices)
