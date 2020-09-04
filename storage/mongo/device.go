package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateDevice implements api.StorageProvider.CreateDevice
func (mdp *MongoStorageProvider) CreateDevice(
	ctx context.Context,
	parent string,
	newDevice *api.Device) (device *api.Device, err error) {

	l, tokens := splitTokens(parent)
	glog.Infof("length=%d tokens=%v", l, tokens)
	if l < 2 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("CreateDevice: %v, 2 required", tokens))
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
		err = info.Inform(err, ErrMongoInsert,
			fmt.Sprintf("CreateDevice: %v", newDevice.Label))
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
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("DeleteDevice: %v, 4 required", tokens))
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
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("GetDevice: 4 required %v", tokens))
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
		err = info.Inform(err, ErrFind, "ListDevices")
		return
	}

	err = cursor.All(ctx, &devices)
	if err != nil {
		err = info.Inform(err, ErrDecode, "ListDevices")
	}
	return
}

// UpdateDevice implements api.StorageProvider.UpdateDevice
func (mdp *MongoStorageProvider) UpdateDevice(ctx context.Context, name string, patch *api.Device) (device *api.Device, err error) {
	device = patch
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, patch)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("UpdateDevice: 4 required %v", tokens))
		return
	}
	filter := bson.D{
		{Key: "domain", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	res, err := mdp.collection.ReplaceOne(ctx, filter, patch)
	if err != nil {
		err = info.Inform(err, ErrDecode,
			fmt.Sprintf("UpdateDevice: %v", name))
		return
	}
	if res.MatchedCount == 0 {
		err = info.Inform(err, ErrNoMatch,
			fmt.Sprintf("UpdateDevice: %v", name))
		return
	}
	return
}