package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateSketch implements api.StorageProvider.CreateSketch
func (mdp *MongoStorageProvider) CreateSketch(
	ctx context.Context,
	parent string,
	newSketch *api.Sketch) (sketch *api.Sketch, err error) {

	l, tokens := splitTokens(parent)
	glog.Infof("length=%d tokens=%v", l, tokens)
	if l < 2 {
		err = info.Inform(err, ErrNotEnough, fmt.Sprintf("%v, 2 required", tokens))
		return
	}

	newSketch.Toolkit = tokens[1]
	b, err := bson.Marshal(newSketch)
	if err != nil {
		err = info.Inform(err, ErrBsonMarshall, newSketch.Label)
		return
	}

	_, err = mdp.Collections[sketchCollectionName].InsertOne(ctx, b)
	if err != nil {
		err = info.Inform(err, ErrMongoInsert,
			fmt.Sprintf("CreateSketch: %v", newSketch.Label))
		return
	}

	sketch = newSketch
	return
}

// DeleteSketch implements api.StorageProvider.DeleteSketch
func (mdp *MongoStorageProvider) DeleteSketch(ctx context.Context, name string) (err error) {
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, name)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("DeleteSketch: %v, 4 required", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	mdp.Collections[sketchCollectionName].FindOneAndDelete(ctx, filter)
	return
}

// GetSketch implements api.StorageProvider.GetSketch
func (mdp *MongoStorageProvider) GetSketch(ctx context.Context, name string) (sketch *api.Sketch, err error) {
	l, tokens := splitTokens(name)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("GetSketch: 4 required %v", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	sketch = &api.Sketch{}
	err = mdp.Collections[sketchCollectionName].FindOne(ctx, filter).Decode(sketch)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetSketch: %v", name))
		return
	}

	return
}

// ListSketches implements api.StorageProvider.ListSketches
func (mdp *MongoStorageProvider) ListSketches(ctx context.Context, parent string) (sketches []*api.Sketch, err error) {
	l, tokens := splitTokens(parent)
	glog.Info(tokens, l)
	filter := bson.D{}
	if l > 1 {
		filter = bson.D{{Key: "model", Value: tokens[1]}}
	}

	sketches = make([]*api.Sketch, 0)
	cursor, err := mdp.Collections[sketchCollectionName].Find(ctx, filter)
	if err != nil {
		err = info.Inform(err, ErrFind, "ListSketches")
		return
	}

	err = cursor.All(ctx, &sketches)
	if err != nil {
		err = info.Inform(err, ErrDecode, "ListSketches")
	}
	return
}

// UpdateSketch implements api.StorageProvider.UpdateSketch
func (mdp *MongoStorageProvider) UpdateSketch(ctx context.Context, name string, patch *api.Sketch) (sketch *api.Sketch, err error) {
	sketch = patch
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, patch)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("UpdateSketch: 4 required %v", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	res, err := mdp.Collections[sketchCollectionName].ReplaceOne(ctx, filter, patch)
	if err != nil {
		err = info.Inform(err, ErrDecode,
			fmt.Sprintf("UpdateSketch: %v", name))
		return
	}
	if res.MatchedCount == 0 {
		err = info.Inform(err, ErrNoMatch,
			fmt.Sprintf("UpdateSketch: %v", name))
		return
	}
	return
}
