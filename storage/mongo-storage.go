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
	ErrMongoClient = errors.New("failed to create client")
	// ErrMongoConnect identifies an error creating mongo client
	ErrMongoConnect = errors.New("failed to connect")
	// ErrDisconnect indentifies a disconnection error
	ErrDisconnect = errors.New("failed to disconnect")
)

// MongoStorageProviderNew creates and returns an instance of MongoStorageProvider
func MongoStorageProviderNew(uri, databaseName string) (mdp *MongoStorageProvider, err error) {
	mdp = &MongoStorageProvider{}
	mdp.client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		err = info.Inform(err, ErrMongoClient, fmt.Sprintf("%v", uri))
		return
	}

	ctx := context.Background()
	err = mdp.client.Connect(ctx)
	if err != nil {
		err = info.Inform(err, ErrMongoConnect, fmt.Sprintf("%v", uri))
		return
	}
	defer mdp.disconnect(ctx)
	mdp.collection = mdp.client.Database(databaseName).Collection("devices")
	return
}

func (mdp *MongoStorageProvider) disconnect(ctx context.Context) {
	if err := mdp.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

var (
	// ErrMongoInsert identifies an error inserting record
	ErrMongoInsert = errors.New("InsertOne")
	// ErrBsonMarshall -
	ErrBsonMarshall = errors.New("json.Marshal")
)

// CreateDevice implements api.StorageProvider.CreateDevice
func (mdp *MongoStorageProvider) CreateDevice(parent, label string, newDevice *api.Device) (device *api.Device, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// encode api.Device to JSON
	newDevice.Domain = parent
	newDevice.Label = label
	b, err := bson.Marshal(newDevice)
	if err != nil {
		err = info.Inform(err, ErrBsonMarshall, newDevice.Label)
	}

	res, err := mdp.collection.InsertOne(ctx, b)
	if err != nil {
		err = info.Inform(err, ErrMongoClient, fmt.Sprintf("%v", newDevice.Label))
	}
	id := res.InsertedID
	fmt.Println(id)
	return
}

// DeleteDevice implements api.StorageProvider.DeleteDevice
func (mdp *MongoStorageProvider) DeleteDevice(parent, label string) (err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// mdp.collection.DeleteOne(ctx, filter interface{})
	// cursor, err := mdp.collection.Find(ctx, bson.D{})
	// if err != nil {
	// 	err = info.Inform(err, ErrFind, "collection.Find")
	// 	return
	// }

	// devices = make([]*api.Device, 0)
	// err = cursor.All(ctx, &devices)
	// if err != nil {
	// 	err = info.Inform(err, ErrDecode, "cursor.All")
	// }
	return
}

// GetDevice implements api.StorageProvider.GetDevice
func (mdp *MongoStorageProvider) GetDevice(name string) (device *api.Device, err error) {
	return
}

// ErrFind -
var ErrFind = errors.New("failed to find")

// ErrDecode -
var ErrDecode = errors.New("failed to decode")

// ListDevices implements api.StorageProvider.ListDevices
func (mdp *MongoStorageProvider) ListDevices(parent string) (devices []*api.Device, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := mdp.collection.Find(ctx, bson.D{})
	if err != nil {
		err = info.Inform(err, ErrFind, "collection.Find")
		return
	}

	devices = make([]*api.Device, 0)
	err = cursor.All(ctx, &devices)
	if err != nil {
		err = info.Inform(err, ErrDecode, "cursor.All")
	}
	return
}

// UpdateDevice implements api.StorageProvider.UpdateDevice
func (mdp *MongoStorageProvider) UpdateDevice(parent string, patch *api.Device) (device *api.Device, err error) {
	return
}

// Authenticate against database
func (mdp *MongoStorageProvider) Authenticate(user, pass, name string, patch *api.Device) (err error) {
	// err = mdp.collection.FindOne(context.TODO(), bson.D{{"username", user}}).Decode(&result)
	return
}
