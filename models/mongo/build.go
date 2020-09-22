package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uriFormat = "mongodb://testing:test@localhost:27017/?authSource=sketchit-test-02"

var deputyVersion = "1.01.00"
var deputyName = "Andy"
var deputyOnly = true
var deputyCreate = false
var dbName = "sketchit-test-03"

func init() {
	flag.StringVar(&dbName, "db", dbName, "database name eg: db=sketchit-test-03")
	flag.StringVar(&deputyName, "name", deputyName, "deputy name eg: name=Andy")
	flag.StringVar(&deputyVersion, "ver", deputyVersion, "deputy version eg: ver=1.01.00")
	flag.BoolVar(&deputyCreate, "create", deputyCreate, "create deputy: create=true, update deputy: create=false")
	flag.BoolVar(&deputyOnly, "only", deputyOnly, "only create/update deputy eg: only=false")
}

func main() {
	// required for glog
	flag.Parse()
	err := build(uriFormat, dbName, deputyName, deputyVersion, deputyOnly, deputyCreate)
	if err != nil {
		glog.Fatalf("failed to build %v on\n\t %v\n", dbName, err)
	}
	glog.Infof("succcessfully built %v on\n\t %v\n", dbName, err)
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

func build(dbURI, dbName, deputyName, deputyVersion string, deputyOnly, deputyCreate bool) (err error) {
	var (
		client *mongo.Client
		db     *mongo.Database
	)

	client, err = mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		err = info.Inform(err, ErrClient, dbURI)
		return
	}
	db = client.Database(dbName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

	var create = func() (err error) {
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

		deputyValidator := bson.M{
			"$jsonSchema": deputySchema,
		}
		_, err = createCollection(ctx, db, dbURI, "deputies", deputyValidator, deputyIndeces)
		if err != nil {
			return
		}
		return
	}

	// create the data base
	if !deputyOnly {
		err = create()
		if err != nil {
			return
		}
	}

	// create a deputy
	collection := db.Collection("deputies")
	deputy := makeDeputy(deputyName, deputyVersion)
	if deputyCreate {
		var res *mongo.InsertOneResult
		res, err = collection.InsertOne(ctx, deputy)
		glog.Info(res, err)
		return
	}
	filter := bson.D{
		{Key: "label", Value: deputy.Label},
		{Key: "version", Value: deputy.Version}}
	var res *mongo.SingleResult
	res = collection.FindOneAndReplace(ctx, filter, deputy)
	err = res.Err()
	glog.Warning(err)
	return
}

func createCollection(
	ctx context.Context,
	db *mongo.Database,
	dbURI string,
	collectionName string,
	validator bson.M,
	indexes []mongo.IndexModel) (collection *mongo.Collection, err error) {

	var opts *options.CreateCollectionOptions = options.CreateCollection()
	opts.SetValidator(validator)
	// opts.SetValidationAction("warn")
	// opts.SetValidationLevel(level string)
	err = db.CreateCollection(ctx, collectionName, opts)
	if err != nil {
		err = info.Inform(err, ErrCreate, fmt.Sprintf("%v %v\n", collectionName, dbURI))
		return
	}

	collection = db.Collection(collectionName)

	names, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		err = info.Inform(err, ErrCreateIndexes, fmt.Sprint(collectionName, dbURI))
		return
	}
	glog.Info(names)

	cur, err := collection.Indexes().List(ctx)
	var results = make([]string, 0)
	cur.All(ctx, results)
	glog.Info(results)
	return
}

func makeDeputy(deputyName, deputyVersion string) (deputy *api.Deputy) {
	deputy = &api.Deputy{
		Label:    deputyName,
		Version:  deputyVersion,
		Skills:   skills,
		Features: features,
	}
	return
}
