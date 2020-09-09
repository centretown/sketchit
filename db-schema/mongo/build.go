package main

import (
	"context"
	"errors"
	"flag"
	"time"

	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test-02"

// var dbname = "devices"

func main() {
	// required for glog
	flag.Parse()
	err := build(testURI)
	if err != nil {
		glog.Fatalf("failed to build %v on\n\t %v\n", testURI, err)
	}
	glog.Infof("succcessfully built %v on\n\t %v\n", testURI, err)
}

var (
	// ErrClient -
	ErrClient = errors.New("failed to create new db client")
	// ErrConnect -
	ErrConnect = errors.New("failed to connect")
	// ErrCreate -
	ErrCreate = errors.New("failed to create validator")
	// ErrCreateIndexes -
	ErrCreateIndexes = errors.New("failed to create indexes")
)

func build(dbURI string) (err error) {
	var (
		client *mongo.Client
		name   = "sketchit-test-02"
		// collection *mongo.Collection
		// b          []byte
		db *mongo.Database
	)

	client, err = mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		err = info.Inform(err, ErrClient, dbURI)
		return
	}

	db = client.Database(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		err = info.Inform(err, ErrConnect, dbURI)
		return
	}

	var disconnect = func() {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}
	defer disconnect()

	deviceValidator := bson.M{
		"$jsonSchema": deviceSchema,
	}

	_, err = createCollection(ctx, db, dbURI, "devices", deviceValidator, deviceIndeces)
	if err != nil {
		return
	}

	sketchValidator := bson.M{
		"$jsonSchema": sketchSchema,
	}
	_, err = createCollection(ctx, db, dbURI, "sketches", sketchValidator, sketchIndeces)
	if err != nil {
		return
	}

	return
}

func createCollection(
	ctx context.Context,
	db *mongo.Database,
	dbURI string,
	dbName string,
	validator bson.M,
	indexes []mongo.IndexModel) (collection *mongo.Collection, err error) {

	opts := options.CreateCollection().SetValidator(validator)
	err = db.CreateCollection(ctx, dbName, opts)
	if err != nil {
		err = info.Inform(err, ErrCreate, dbURI)
		return
	}

	collection = db.Collection(dbName)

	names, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		err = info.Inform(err, ErrCreateIndexes, dbName)
		return
	}
	glog.Info(names)

	cur, err := collection.Indexes().List(ctx)
	var results = make([]string, 0)
	cur.All(ctx, results)
	glog.Info(results)
	return
}
