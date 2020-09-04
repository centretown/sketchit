package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
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
