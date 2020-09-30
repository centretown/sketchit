package storage

import (
	"context"
	"fmt"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/info"
)

// GetDeputy implements api.StorageProvider.GetDeputy
func (mdp *MongoStorageProvider) GetDeputy(ctx context.Context, name string) (deputy *api.Deputy, err error) {
	deputy = &api.Deputy{}

	collector, ok := mdp.Collectors[deputiesName]
	if !ok {
		err = info.Inform(err, ErrCollection, fmt.Sprintf("List %v", deputiesName))
		return
	}

	filter := collector.Filter(name)
	// filter := bson.D{}
	res := collector.Collection.FindOne(ctx, filter)
	if res == nil {
		err = info.Inform(err, ErrFind, fmt.Sprintf("GetDeputy FindOne nil response: %v", name))
		return
	}

	err = res.Err()
	if err != nil {
		err = info.Inform(err, ErrFind, fmt.Sprintf("GetDeputy FindOne: %v", name))
		return
	}

	err = res.Decode(deputy)
	if err != nil {
		err = info.Inform(err, ErrDecode, fmt.Sprintf("GetDeputy Decode: %v", name))
		return
	}

	return
}
