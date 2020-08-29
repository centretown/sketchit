package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

// MongoStorageProvider implements the StorageProvider interface
type MongoStorageProvider struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	// ErrMongoClient identifies an error creating mongo client
	ErrMongoClient error = errors.New("from NewClient, ApplyURI")
)

// MongoStorageProviderNew creates and returns an instance of MongoStorageProvider
func MongoStorageProviderNew(uri, name string) (provider api.StorageProvider, err error) {
	mdp := &MongoStorageProvider{}
	mdp.client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		err = info.Inform(err, ErrMongoClient,
			fmt.Sprintf("%v", uri))
		return
	}
	mdp.collection = mdp.client.Database(name).Collection("devices")
	provider = mdp
	return
}

var (
	// ErrMongoInsert identifies an error inserting record
	ErrMongoInsert error = errors.New("from InsertOne")
)

// CreateDevice implements api.StorageProvider.CreateDevice
func (mdp *MongoStorageProvider) CreateDevice(parent, label string, newDevice *api.Device) (device *api.Device, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := mdp.collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		err = info.Inform(err, ErrMongoClient,
			fmt.Sprintf("%v", label))
		return
	}
	id := res.InsertedID
	fmt.Println(id)
	return
}

// DeleteDevice implements api.StorageProvider.DeleteDevice
func (mdp *MongoStorageProvider) DeleteDevice(parent, label string) (err error) {
	return err
}

// GetDevice implements api.StorageProvider.GetDevice
func (mdp *MongoStorageProvider) GetDevice(name string) (device *api.Device, err error) {
	return
}

// ListDevices implements api.StorageProvider.ListDevices
func (mdp *MongoStorageProvider) ListDevices(parent string) (devices []*api.Device, err error) {
	return
}

// UpdateDevice implements api.StorageProvider.UpdateDevice
func (mdp *MongoStorageProvider) UpdateDevice(parent string, patch *api.Device) (device *api.Device, err error) {
	return
}
