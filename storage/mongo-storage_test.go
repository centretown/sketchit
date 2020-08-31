package storage

import (
	"context"
	"testing"
	"time"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

func TestMongoStorageProvider(t *testing.T) {
	databaseName := "sketchit-test"
	mdp, err := MongoStorageProviderNew(testURI, databaseName)
	if err != nil {
		t.Fatalf("MongoStorageProviderNew: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = mdp.client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mdp.disconnect(ctx)

	// err = mdp.CreateModel(ctx, databaseName)
	// if err != nil {
	// 	t.Fatal(err)
	// }

}

// devices, err := mdp.ListDevices("")
// if err != nil {
// 	t.Fatalf("ListDevices: %v", err)
// }

// for _, d := range devices {
// 	t.Logf("label: %v", d.Label)
// }
