package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDeputy implements api.StorageProvider.GetDeputy
func (mdp *MongoStorageProvider) GetDeputy(ctx context.Context, name string) (deputy *api.Deputy, err error) {
	filter := bson.D{{Key: "label", Value: name}}
	deputy = &api.Deputy{}
	res := mdp.Collections[deputyCollectionName].FindOne(ctx, filter)
	if res == nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetDeputy FindOne: %v", "nil response"))
		return
	}

	err = res.Err()
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetDeputy FindOne: %v", name))
		return
	}

	err = res.Decode(deputy)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetDeputy Decode: %v", name))
		return
	}

	return
}
