package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateProcess implements api.StorageProvider.CreateProcess
func (mdp *MongoStorageProvider) CreateProcess(
	ctx context.Context,
	parent string,
	newProcess *api.Process) (process *api.Process, err error) {

	l, tokens := splitTokens(parent)
	glog.Infof("length=%d tokens=%v", l, tokens)
	if l < 2 {
		err = info.Inform(err, ErrNotEnough, fmt.Sprintf("%v, 2 required", tokens))
		return
	}

	newProcess.Model = tokens[1]
	b, err := bson.Marshal(newProcess)
	if err != nil {
		err = info.Inform(err, ErrBsonMarshall, newProcess.Label)
		return
	}

	_, err = mdp.collection.InsertOne(ctx, b)
	if err != nil {
		err = info.Inform(err, ErrMongoInsert,
			fmt.Sprintf("CreateProcess: %v", newProcess.Label))
		return
	}

	process = newProcess
	return
}

// DeleteProcess implements api.StorageProvider.DeleteProcess
func (mdp *MongoStorageProvider) DeleteProcess(ctx context.Context, name string) (err error) {
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, name)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("DeleteProcess: %v, 4 required", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	mdp.collection.FindOneAndDelete(ctx, filter)
	return
}

// GetProcess implements api.StorageProvider.GetProcess
func (mdp *MongoStorageProvider) GetProcess(ctx context.Context, name string) (process *api.Process, err error) {
	l, tokens := splitTokens(name)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("GetProcess: 4 required %v", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	process = &api.Process{}
	err = mdp.collection.FindOne(ctx, filter).Decode(process)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetProcess: %v", name))
		return
	}

	return
}

// ListProcesses implements api.StorageProvider.ListProcesses
func (mdp *MongoStorageProvider) ListProcesses(ctx context.Context, parent string) (processes []*api.Process, err error) {
	l, tokens := splitTokens(parent)
	glog.Info(tokens, l)
	filter := bson.D{}
	if l > 1 {
		filter = bson.D{{Key: "model", Value: tokens[1]}}
	}

	processes = make([]*api.Process, 0)
	cursor, err := mdp.collection.Find(ctx, filter)
	if err != nil {
		err = info.Inform(err, ErrFind, "ListProcesses")
		return
	}

	err = cursor.All(ctx, &processes)
	if err != nil {
		err = info.Inform(err, ErrDecode, "ListProcesses")
	}
	return
}

// UpdateProcess implements api.StorageProvider.UpdateProcess
func (mdp *MongoStorageProvider) UpdateProcess(ctx context.Context, name string, patch *api.Process) (process *api.Process, err error) {
	process = patch
	l, tokens := splitTokens(name)
	glog.Info(tokens, l, patch)
	if l < 4 {
		err = info.Inform(err, ErrNotEnough,
			fmt.Sprintf("UpdateProcess: 4 required %v", tokens))
		return
	}
	filter := bson.D{
		{Key: "model", Value: tokens[1]},
		{Key: "label", Value: tokens[3]}}

	res, err := mdp.collection.ReplaceOne(ctx, filter, patch)
	if err != nil {
		err = info.Inform(err, ErrDecode,
			fmt.Sprintf("UpdateProcess: %v", name))
		return
	}
	if res.MatchedCount == 0 {
		err = info.Inform(err, ErrNoMatch,
			fmt.Sprintf("UpdateProcess: %v", name))
		return
	}
	return
}
