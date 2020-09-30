package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

// MongoStorageProvider implements the StorageProvider interface
type MongoStorageProvider struct {
	client     *mongo.Client
	Name       string
	URI        string
	Collectors map[string]*Collector
}

const (
	devicesName  = "devices"
	sketchesName = "sketches"
	deputiesName = "deputies"
)

// MongoStorageProviderNew creates and returns an instance of MongoStorageProvider
func MongoStorageProviderNew(uri, databaseName, authSource string) (mdp *MongoStorageProvider, err error) {
	mdp = &MongoStorageProvider{}
	mdp.Name = databaseName
	mdp.URI = uri
	mdp.client, err = mongo.NewClient(options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthSource: authSource, Username: "testing", Password: "test"}))
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

	// map supported collections
	mdp.Collectors = make(map[string]*Collector)

	mdp.Collectors[devicesName] = &Collector{
		Collection: mdp.client.Database(mdp.Name).Collection(devicesName),
		NewItem:    func() protoreflect.ProtoMessage { return &api.Device{} },
		Filter: func(parent string) bson.D {
			return makeFilter(parent, "sector", "label")
		},
	}
	mdp.Collectors[sketchesName] = &Collector{
		Collection: mdp.client.Database(mdp.Name).Collection(sketchesName),
		NewItem:    func() protoreflect.ProtoMessage { return &api.Sketch{} },
		Filter: func(parent string) bson.D {
			return makeFilter(parent, "toolkit", "label")
		},
	}
	mdp.Collectors[deputiesName] = &Collector{
		Collection: mdp.client.Database(mdp.Name).Collection(deputiesName),
		NewItem:    func() protoreflect.ProtoMessage { return &api.Deputy{} },
		Filter: func(name string) (filter bson.D) {
			return bson.D{{Key: "label", Value: name}}
		},
	}
	return
}

// ListCollections list the collections in the current database
func (mdp *MongoStorageProvider) ListCollections(ctx context.Context, name string) (collections []*api.Collection, err error) {
	db := mdp.client.Database(mdp.Name)
	filter := bson.D{}
	opts := &options.ListCollectionsOptions{}
	cursor, err := db.ListCollections(ctx, filter, opts)
	if err != nil {
		err = info.Inform(err, ErrCollectionNames, "ListCollections")
		return
	}

	for cursor.Next(ctx) {
		c := &MongoCollection{}
		cursor.Decode(c)
		collections = append(collections, c.MongoCollectionNew())
	}
	return
}

func (mdp *MongoStorageProvider) getCollector(parent string) (collector *Collector, ok bool) {
	tokens, length := splitParent(parent)
	if length < CollectionName+1 {
		return
	}
	collector, ok = mdp.Collectors[tokens[CollectionName]]
	return
}

// Authenticate against database
func (mdp *MongoStorageProvider) Authenticate(user, pass, name string, patch *api.Device) (err error) {
	// err = mdp.collection.FindOne(context.TODO(), bson.D{{"username", user}}).Decode(&result)
	return
}

// SplitParent -
func splitParent(parent string) (tokens []string, length int) {
	sep := "/"
	parent = strings.TrimRight(parent, sep)
	tokens = strings.Split(parent, sep)
	length = len(tokens)
	return
}

func makeFilter(parent, parentName, labelName string) (filter bson.D) {
	tokens, l := splitParent(parent)
	if l > 1 {
		filter = bson.D{{Key: parentName, Value: tokens[1]}}
		if l > 3 {
			filter = append(filter, bson.E{Key: labelName, Value: tokens[3]})
		}
	}
	return
}
