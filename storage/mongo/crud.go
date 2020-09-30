package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/info"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/anypb"
)

// Create implements api.StorageProvider.Create
func (mdp *MongoStorageProvider) Create(ctx context.Context, parent string, newAny *anypb.Any) (item *anypb.Any, err error) {

	collector, ok := mdp.getCollector(parent)
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("Create %v", parent))
		return
	}

	newItem := collector.NewItem()
	newAny.UnmarshalTo(newItem)
	b, err := bson.Marshal(newItem)
	if err != nil {
		err = info.Inform(err, ErrBsonMarshall, parent)
		return
	}

	_, err = collector.Collection.InsertOne(ctx, b)
	if err != nil {
		err = info.Inform(err, ErrMongoInsert, fmt.Sprintf("Create: %v", parent))
		return
	}
	item = newAny
	return
}

// Get implements api.StorageProvider.Get
func (mdp *MongoStorageProvider) Get(ctx context.Context, name string) (item *anypb.Any, err error) {
	collector, ok := mdp.getCollector(name)
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("Get %v", name))
		return
	}

	filter := collector.Filter(name)
	kind := collector.NewItem()
	var res *mongo.SingleResult = collector.Collection.FindOne(ctx, filter)
	err = res.Err()
	if err != nil {
		err = info.Inform(err, ErrFind, fmt.Sprintf("Get %v", name))
		return
	}
	err = res.Decode(kind)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("Get %v", name))
		return
	}
	item = &anypb.Any{}
	err = item.MarshalFrom(kind)
	if err != nil {
		err = info.Inform(err, ErrDecodeAny, fmt.Sprintf("Get %v", name))
		return
	}
	return
}

// List implements api.StorageProvider.List
func (mdp *MongoStorageProvider) List(ctx context.Context, parent string) (list []*anypb.Any, err error) {

	collector, ok := mdp.getCollector(parent)
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("List %v", parent))
		return
	}

	filter := collector.Filter(parent)
	cursor, err := collector.Collection.Find(ctx, filter)
	if err != nil {
		err = info.Inform(err, ErrFind, fmt.Sprintf("List %v", parent))
		return
	}
	defer cursor.Close(ctx)

	list = make([]*anypb.Any, 0)
	for cursor.Next(ctx) {
		kind := collector.NewItem()
		err = cursor.Decode(kind)
		if err != nil {
			err = info.Inform(err, ErrDecode, fmt.Sprintf("List %v", parent))
			return
		}
		any := &anypb.Any{}
		err = any.MarshalFrom(kind)
		if err != nil {
			err = info.Inform(err, ErrDecodeAny, fmt.Sprintf("List %v", parent))
			return
		}
		list = append(list, any)
	}

	return
}

// Update implements api.StorageProvider.Update
func (mdp *MongoStorageProvider) Update(ctx context.Context, name string, patch *anypb.Any) (item *anypb.Any, err error) {
	collector, ok := mdp.getCollector(name)
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("List %v", name))
		return
	}

	filter := collector.Filter(name)
	kind := collector.NewItem()

	err = patch.UnmarshalTo(kind)
	if err != nil {
		err = info.Inform(err, ErrDecodeAny, fmt.Sprintf("Update: %v", name))
		return
	}

	res, err := collector.Collection.ReplaceOne(ctx, filter, kind)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("Update: %v", name))
		return
	}
	if res.MatchedCount == 0 {
		err = info.Inform(err, ErrNoMatch, fmt.Sprintf("Update: %v", name))
		return
	}
	item = patch
	return
}

// Delete implements api.StorageProvider.Delete
func (mdp *MongoStorageProvider) Delete(ctx context.Context, name string) (err error) {
	collector, ok := mdp.getCollector(name)
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("Delete %v", name))
		return
	}

	filter := collector.Filter(name)
	res := collector.Collection.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		err = info.Inform(res.Err(), ErrNoMatch, fmt.Sprintf("Delete: %v", name))
		return
	}
	return
}
