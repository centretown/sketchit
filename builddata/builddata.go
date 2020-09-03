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

var testURI = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test"
var dbname = "devices"

func main() {
	// required for glog
	flag.Parse()
	err := build(dbname, testURI)
	if err != nil {
		glog.Fatalf("failed to build %v on\n\t %v: %v\n", dbname, testURI, err)
	}
	glog.Infof("succcessfully built %v on\n\t %v: %v\n", dbname, testURI, err)
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

func build(dbName, dbURI string) (err error) {
	type indexTrait struct {
		idUnique   bool
		pathUnique bool
		pathName   string
	}

	var (
		client *mongo.Client
		name   = "sketchit-test"
		// collection *mongo.Collection
		// b          []byte
		db          *mongo.Database
		deviceTrait = &indexTrait{
			idUnique:   true,
			pathUnique: true,
			pathName:   "pathIndex",
		}
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

	//
	var disconnect = func() {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}
	defer disconnect()

	deviceSchema := bson.M{
		"bsonType": "object",
		"required": []string{"domain", "label", "model"},
		"properties": bson.M{
			"domain": bson.M{
				"bsonType":    "string",
				"description": "the name of the domain, which is required and must be a string",
			},
			"label": bson.M{
				"bsonType":    "string",
				"description": "the name of the device, which is required and must be a string",
			},
			"model": bson.M{
				"bsonType":    "string",
				"description": "the device model type, which is required and must be a string",
			},
			"pins": bson.M{
				"bsonType":    "array",
				"description": "the managed pins on this device",
				"items": bson.M{
					"title":       "Pin",
					"bsonType":    "object",
					"description": "the pin definition",
					"required":    []string{"id", "label"},
					"uniqueItems": true,
					"properties": bson.M{
						"id": bson.M{
							"bsonType":    "number",
							"description": "the pin number, which is required and must be a number",
						},
						"label": bson.M{
							"bsonType":    "string",
							"description": "the label assigned to the pin, which is required and must be a string",
						},
						"purpose": bson.M{
							"bsonType":    "string",
							"description": "the device model type, which is required and must be a string",
						},
					},
				},
			},
		},
	}
	validator := bson.M{
		"$jsonSchema": deviceSchema,
	}
	opts := options.CreateCollection().SetValidator(validator)

	if err = db.CreateCollection(ctx, dbName, opts); err != nil {
		if err != nil {
			err = info.Inform(err, ErrCreate, dbURI)
			return
		}
	}

	var deviceIndeces = []mongo.IndexModel{
		{
			Keys: bson.M{
				"domain": 1,
				"label":  1,
			},
			Options: &options.IndexOptions{
				Name:   &deviceTrait.pathName,
				Unique: &deviceTrait.pathUnique,
			},
		},
		{
			Keys: bson.M{
				"model": 1,
			},
		},
	}

	collection := db.Collection(dbName)
	names, err := collection.Indexes().CreateMany(ctx, deviceIndeces)
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
