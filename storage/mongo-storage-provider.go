package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

// MongoStorageProvider implements the StorageProvider interface
type MongoStorageProvider struct {
	client     *mongo.Client
	Name       string
	URI        string
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
	mdp.Name = databaseName
	mdp.URI = uri
	mdp.client, err = mongo.NewClient(options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthSource: databaseName, Username: "testing", Password: "test"}))
	if err != nil {
		err = info.Inform(err, ErrMongoClient, fmt.Sprintf("%v: %v", databaseName, uri))
		return
	}
	ctx := context.Background()
	err = mdp.client.Connect(ctx)
	if err != nil {
		err = info.Inform(err, ErrMongoConnect, "failed to connect client")
		return
	}
	mdp.collection = mdp.client.Database(mdp.Name).Collection("devices")
	return
}

var (
	// ErrMongoInsert identifies an error inserting record
	ErrMongoInsert = errors.New("failed to insert")
	// ErrBsonMarshall -
	ErrBsonMarshall = errors.New("failed to create json")
	// ErrCollection -
	ErrCollection = errors.New("failed to get collection")
)

// CreateDevice implements api.StorageProvider.CreateDevice
func (mdp *MongoStorageProvider) CreateDevice(
	ctx context.Context,
	parent string,
	newDevice *api.Device) (device *api.Device, err error) {

	l, tokens := splitTokens(parent)
	glog.Infof("length=%d tokens=%v", l, tokens)
	if l < 2 {
		err = info.Inform(err, ErrFind,
			fmt.Sprintf("update device: too few arguments %v, 2 required", tokens))
		return
	}

	newDevice.Domain = tokens[1]
	b, err := bson.Marshal(newDevice)
	if err != nil {
		err = info.Inform(err, ErrBsonMarshall, newDevice.Label)
		return
	}

	_, err = mdp.collection.InsertOne(ctx, b)
	if err != nil {
		err = info.Inform(err, ErrMongoInsert, fmt.Sprintf("%v", newDevice.Label))
		return
	}

	device = newDevice
	return
}

// DeleteDevice implements api.StorageProvider.DeleteDevice
func (mdp *MongoStorageProvider) DeleteDevice(ctx context.Context, name string) (err error) {
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, name)
	if l < 4 {
		err = info.Inform(err, ErrFind,
			fmt.Sprintf("update device: too few arguments %v, 4 required", tokens))
		return
	}
	filter := bson.D{
		{Key: "domain", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	mdp.collection.FindOneAndDelete(ctx, filter)
	return
}

// GetDevice implements api.StorageProvider.GetDevice
func (mdp *MongoStorageProvider) GetDevice(ctx context.Context, name string) (device *api.Device, err error) {
	l, tokens := splitTokens(name)
	if l < 4 {
		err = info.Inform(err, ErrFind, "not enough arguments")
		return
	}
	filter := bson.D{
		{Key: "domain", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	device = &api.Device{}
	err = mdp.collection.FindOne(ctx, filter).Decode(device)
	if err != nil {
		err = info.Inform(err, ErrDecode, "get device")
		return
	}

	return
}

var (
	// ErrFind -
	ErrFind = errors.New("failed to find")
	// ErrDecode -
	ErrDecode = errors.New("failed to decode")
)

// ListDevices implements api.StorageProvider.ListDevices
func (mdp *MongoStorageProvider) ListDevices(ctx context.Context, parent string) (devices []*api.Device, err error) {
	l, tokens := splitTokens(parent)
	glog.Info(tokens, l)
	filter := bson.D{}
	if l > 1 {
		filter = bson.D{{Key: "domain", Value: tokens[1]}}
	}

	devices = make([]*api.Device, 0)
	cursor, err := mdp.collection.Find(ctx, filter)
	if err != nil {
		err = info.Inform(err, ErrFind, "devices")
		return
	}

	err = cursor.All(ctx, &devices)
	if err != nil {
		err = info.Inform(err, ErrDecode, "all devices")
	}
	return
}

// ErrNoMatch -
var ErrNoMatch = errors.New("No records matched")

// UpdateDevice implements api.StorageProvider.UpdateDevice
func (mdp *MongoStorageProvider) UpdateDevice(ctx context.Context, name string, patch *api.Device) (device *api.Device, err error) {
	device = patch
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, patch)
	if l < 4 {
		err = info.Inform(err, ErrFind,
			fmt.Sprintf("update device: too few arguments %v, 4 required", tokens))
		return
	}
	filter := bson.D{
		{Key: "domain", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	res, err := mdp.collection.ReplaceOne(ctx, filter, patch)
	if err != nil {
		err = info.Inform(err, ErrDecode,
			fmt.Sprintf("update device: %v", tokens))
		return
	}
	if res.MatchedCount == 0 {
		err = info.Inform(err, ErrNoMatch,
			fmt.Sprintf("update device: %v", patch))
		return
	}
	return
}

// Authenticate against database
func (mdp *MongoStorageProvider) Authenticate(user, pass, name string, patch *api.Device) (err error) {
	// err = mdp.collection.FindOne(context.TODO(), bson.D{{"username", user}}).Decode(&result)
	return
}

func splitTokens(source string) (l int, tokens []string) {
	source = strings.TrimRight(source, "/")
	tokens = strings.Split(source, "/")
	l = len(tokens)
	return
}
