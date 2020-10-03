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
	"google.golang.org/protobuf/types/known/anypb"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test-02"

var newDevice = &api.Device{
	Sector:  "home",
	Label:   "esp32-04",
	Toolkit: "ESP32",
	Ip:      "192.168.1.200:8080",
	Port:    "esp32-01",
	Pins: []*api.Device_Pin{
		{Pin: 2, Label: "LED", Purpose: "activity indicator"},
		{Pin: 5, Label: "TX", Purpose: "soft serial transmitter"},
		{Pin: 6, Label: "RX", Purpose: "soft serial receiver"},
	},
}

var (
	databaseName = "sketchit-test-03"
	authSource   = "sketchit-test-02"
)

func TestMongoGetDeputy(t *testing.T) {
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

func showMongoSchema(sch *Schema, name string, level *indent) {
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
			showMongoSchema(&s, fmt.Sprint(i), level)
		}
		*level--
	}
	*level++
	if sch.Items != nil {
		showMongoSchema(sch.Items, "items", level)
	}
	for k, s := range sch.Properties {
		showMongoSchema(&s, k, level)
	}
	*level--
}

func TestDevices(t *testing.T) {
	mdp, err := MongoStorageProviderNew(testURI, databaseName, authSource)
	if err != nil {
		t.Fatalf("TestDevices: %v", err)
	}

	ctx := context.Background()

	// delete previously added
	name := fmt.Sprintf("sectors/%s/devices/%s", newDevice.Sector, newDevice.Label)
	err = mdp.Delete(ctx, name)
	if err != nil {
		t.Logf("TestDevices delete: %v", err)
	}
	t.Log("DELETE")

	// create
	parent := "sectors/home/devices"
	any := &anypb.Any{}
	any.MarshalFrom(newDevice)
	device, err := mdp.Create(ctx, parent, any)
	if err != nil {
		t.Fatalf("TestDevices create: %v", err)
	}
	t.Log("CREATE", device)

	// list
	list, err := mdp.List(ctx, parent)
	if err != nil {
		t.Fatalf("TestDevices list: %v", err)
	}

	for _, any := range list {
		dev := &api.Device{}
		any.UnmarshalTo(dev)
		t.Log("LIST", dev)
	}

	parent = "sectors/*/devices"
	list, err = mdp.List(ctx, parent)
	if err != nil {
		t.Fatalf("TestDevices list wildcard: %v", err)
	}

	for _, any := range list {
		dev := &api.Device{}
		any.UnmarshalTo(dev)
		t.Log("LIST WILDCARD", dev)
	}

	// get
	item, err := mdp.Get(ctx, name)
	if err != nil {
		t.Fatalf("TestDevices get: %v", err)
	}
	t.Log(item)
	dev := &api.Device{}
	item.UnmarshalTo(dev)
	t.Log("GET", dev)

	// update
	dev.Port = dev.Label
	patch := &anypb.Any{}
	patch.MarshalFrom(dev)
	item, err = mdp.Update(ctx, name, patch)
	if err != nil {
		t.Fatalf("TestDevices update: %v", err)
	}

	item.UnmarshalTo(dev)
	t.Log("UPDATE", dev)
}

func TestSketches(t *testing.T) {
	mdp, err := MongoStorageProviderNew(testURI, databaseName, authSource)
	if err != nil {
		t.Fatalf("TestSketches: %v", err)
	}

	ctx := context.Background()

	// delete previously added
	name := fmt.Sprintf("toolkits/%s/sketches/%s", newSketch2.Toolkit, newSketch2.Label)
	err = mdp.Delete(ctx, name)
	if err != nil {
		t.Logf("TestDevices delete: %v", err)
	}
	t.Log("DELETE")

	name = fmt.Sprintf("toolkits/%s/sketches/%s", newSketch.Toolkit, newSketch.Label)
	err = mdp.Delete(ctx, name)
	if err != nil {
		t.Logf("TestDevices delete: %v", err)
	}
	t.Log("DELETE")

	// create
	parent := "toolkits/ESP32/sketches"
	any := &anypb.Any{}
	any.MarshalFrom(newSketch2)
	any, err = mdp.Create(ctx, parent, any)
	if err != nil {
		t.Fatalf("TestSketches create: %v", err)
	}
	t.Log("CREATE", any)

	any.MarshalFrom(newSketch)
	any, err = mdp.Create(ctx, parent, any)
	if err != nil {
		t.Fatalf("TestSketches create: %v", err)
	}
	t.Log("CREATE", any)

	// list
	parent = "toolkits/ESP32/sketches"
	list, err := mdp.List(ctx, parent)
	if err != nil {
		t.Fatalf("TestSketches list: %v", err)
	}

	for _, any := range list {
		sketch := &api.Sketch{}
		any.UnmarshalTo(sketch)
		t.Log("LIST", sketch)
	}

	// list wildcard
	parent = "toolkits/*/sketches"
	list, err = mdp.List(ctx, parent)
	if err != nil {
		t.Fatalf("TestSketches wildcard list: %v", err)
	}

	for _, any := range list {
		sketch := &api.Sketch{}
		any.UnmarshalTo(sketch)
		t.Log("LIST WILDCARD", sketch)
	}

	// get
	name = fmt.Sprintf("toolkits/%s/sketches/%s", newSketch2.Toolkit, newSketch2.Label)
	item, err := mdp.Get(ctx, name)
	if err != nil {
		t.Fatalf("TestSketches get: %v", err)
	}
	t.Log(item)
	sketch := &api.Sketch{}
	item.UnmarshalTo(sketch)
	t.Log("GET", sketch)

	// update
	sketch.Purpose = "CHANGED OUR MIND"
	patch := &anypb.Any{}
	patch.MarshalFrom(sketch)
	item, err = mdp.Update(ctx, name, patch)
	if err != nil {
		t.Fatalf("TestDevices update: %v", err)
	}

	item.UnmarshalTo(sketch)
	t.Log("UPDATE", sketch)

	// get
	name = fmt.Sprintf("toolkits/%s/sketches/%s", newSketch2.Toolkit, newSketch2.Label)
	item, err = mdp.Get(ctx, name)
	if err != nil {
		t.Fatalf("TestSketches get: %v", err)
	}
	t.Log(item)
	sketch = &api.Sketch{}
	item.UnmarshalTo(sketch)
	t.Log("GET", sketch)
}

var newSketch = &api.Sketch{
	Toolkit: "NANO",
	Label:   "blink02",
	Purpose: "Blinks led at regular intervals",
	Device:  "sectors/home/devices/esp32-04",
	Setup: []*api.Action{{
		Operation: api.Operation_delay,
		Arguments: []int32{500},
	}},
	Loop: []*api.Action{{
		Operation: api.Operation_delay,
		Arguments: []int32{500},
	}},
}

var newSketch2 = &api.Sketch{
	Toolkit: "ESP32",
	Label:   "blink03",
	Purpose: "Blinks led at regular intervals",
	Device:  "sectors/home/devices/esp32-04",
	Setup: []*api.Action{
		{
			Operation: api.Operation_mode,
			Arguments: []int32{2,
				int32(api.Signal_digital),
				int32(api.Mode_output)},
		},
	},
	Loop: []*api.Action{
		{
			Operation: api.Operation_pin,
			Arguments: []int32{2,
				int32(api.Signal_digital),
				int32(api.Mode_output),
				int32(api.Digital_high)},
		},
		{
			Operation: api.Operation_delay,
			Arguments: []int32{500},
		},
		{
			Operation: api.Operation_pin,
			Arguments: []int32{2,
				int32(api.Signal_digital),
				int32(api.Mode_output),
				int32(api.Digital_low)},
		},
		{
			Operation: api.Operation_delay,
			Arguments: []int32{500},
		},
	},
}
