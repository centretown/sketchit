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
)

// "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"

// MongoStorageProvider implements the StorageProvider interface
type MongoStorageProvider struct {
	client      *mongo.Client
	Name        string
	URI         string
	Collections map[string]*mongo.Collection
}

var deviceCollectionName = "devices"
var sketchCollectionName = "sketches"

// MongoStorageProviderNew creates and returns an instance of MongoStorageProvider
func MongoStorageProviderNew(uri, databaseName string) (mdp *MongoStorageProvider, err error) {
	mdp = &MongoStorageProvider{}
	mdp.Name = databaseName
	mdp.URI = uri
	mdp.Collections = make(map[string]*mongo.Collection)
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

	// map supported collections
	mdp.Collections[deviceCollectionName] = mdp.client.Database(mdp.Name).Collection(deviceCollectionName)
	mdp.Collections[sketchCollectionName] = mdp.client.Database(mdp.Name).Collection(sketchCollectionName)
	return
}

type result struct{}

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
	// var colls = &MongoCollection{}
	// err = cursor.All(ctx, colls)
	// if err != nil {
	// err = info.Inform(err, ErrCollectionNames, "ListCollections")
	// return
	// }
	for cursor.Next(ctx) {
		c := &MongoCollection{}
		cursor.Decode(c)
		collections = append(collections, c.MongoCollectionNew())
		//
		// fmt.Printf("Name: %s, Type: %s\n", coll.Name, coll.Type)
		// sch := c.Options.Validator.JSONSchema
		// level := indent(0)
		// showMongoSchema(sch, c.Name, &level)
		// showSchema(coll.Model, &level)

	}
	// glog.Infof("ListCollections %+v", res)
	//names = []string{deviceCollectionName, sketchCollectionName}
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
