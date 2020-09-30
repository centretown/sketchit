package storage

import (
	"fmt"
	"strings"

	"github.com/centretown/sketchit/api"
)

// MongoCollection -
type MongoCollection struct {
	Name    string `bson:"name"`
	Type    string `bson:"type"`
	Options struct {
		Validator struct {
			JSONSchema Schema `bson:"$jsonSchema"`
		} `bson:"validator"`
	} `bson:"options"`
	Info struct {
		ReadOnly bool `bson:"readOnly"`
		// UUID     string `json:"uuid"`
	} `bson:"info"`
	IDIndex struct {
		V   int `bson:"v"`
		Key struct {
			ID int `bson:"_id"`
		} `json:"key"`
		Name string `bson:"name"`
	} `bson:"idIndex"`
}

// MongoCollectionNew convert to message
func (coll *MongoCollection) MongoCollectionNew() (collection *api.Collection) {
	collection = &api.Collection{
		Name:     coll.Name,
		Type:     coll.Type,
		ReadOnly: coll.Info.ReadOnly,
	}
	sch := coll.Options.Validator.JSONSchema
	collection.Model = sch.makeModel(coll.Name)
	return
}

// Schema -
type Schema struct {
	Title       string            `bson:"title"`
	BsonType    string            `bson:"bsonType"`
	Description string            `bson:"description"`
	Required    []string          `bson:"required"`
	Enum        []string          `bson:"enum"`
	OneOf       []Schema          `bson:"oneOf"`
	UniqueItems bool              `bson:"uniqueItems"`
	Properties  map[string]Schema `bson:"properties"`
	Items       *Schema           `bson:"items"`
}

// makeModel converts MongoSchema to api.Model
func (sch *Schema) makeModel(name string) (model *api.Model) {
	model = &api.Model{
		Label:       strings.ToLower(sch.Title),
		Title:       sch.Title,
		Type:        sch.BsonType,
		Description: sch.Description,
		UniqueItems: sch.UniqueItems,
		Required:    sch.Required,
		Options:     sch.Enum,
	}

	if sch.Items != nil {
		model.Items = sch.Items.makeModel("items")
	}

	if len(sch.OneOf) > 0 {
		model.OneOf = make([]*api.Model, 0, len(sch.OneOf))
		for i, s := range sch.OneOf {
			model.OneOf = append(model.OneOf,
				s.makeModel(fmt.Sprintf("option-%0d", i+1)))
		}
	}

	if len(sch.Properties) > 0 {
		model.Properties = make([]*api.Model, len(sch.Required), len(sch.Properties))
		for key, property := range sch.Properties {
			vacant := true
			for i, required := range sch.Required {
				if key == required {
					model.Properties[i] = property.makeModel(key)
					vacant = false
					break
				}
			}
			if vacant {
				model.Properties = append(model.Properties, property.makeModel(key))
			}
		}
	}

	return
}
