package storage

import (
	"fmt"

	"github.com/centretown/sketchit/api"
)

// MongoSchema -
type MongoSchema struct {
	Title       string                 `bson:"title"`
	BsonType    string                 `bson:"bsonType"`
	Description string                 `bson:"description"`
	Required    []string               `bson:"required"`
	Enum        []string               `bson:"enum"`
	OneOf       []MongoSchema          `bson:"oneOf"`
	UniqueItems bool                   `bson:"uniqueItems"`
	Properties  map[string]MongoSchema `bson:"properties"`
	Items       *MongoSchema           `bson:"items"`
}

// MakeSchema -
func (sch *MongoSchema) makeSchema(name string) (schema *api.Schema) {
	schema = &api.Schema{
		Name:        name,
		Title:       sch.Title,
		Type:        sch.BsonType,
		Description: sch.Description,
		UniqueItems: sch.UniqueItems,
		Required:    sch.Required,
		Enum:        sch.Enum,
	}

	if sch.Items != nil {
		schema.Items = sch.Items.makeSchema("items")
	}

	if len(sch.OneOf) > 0 {
		schema.OneOf = make([]*api.Schema, 0, len(sch.OneOf))
		for i, s := range sch.OneOf {
			schema.OneOf = append(schema.OneOf,
				s.makeSchema(fmt.Sprintf("option-%0d", i+1)))
		}
	}

	if len(sch.Properties) > 0 {
		schema.Properties = make([]*api.Schema, len(sch.Required), len(sch.Properties))
		for key, property := range sch.Properties {
			vacant := true
			for i, required := range sch.Required {
				if key == required {
					schema.Properties[i] = property.makeSchema(key)
					vacant = false
					break
				}
			}
			if vacant {
				schema.Properties = append(schema.Properties, property.makeSchema(key))
			}
		}
	}

	return
}

// MongoCollection -
type MongoCollection struct {
	Name    string `bson:"name"`
	Type    string `bson:"type"`
	Options struct {
		Validator struct {
			JSONSchema MongoSchema `bson:"$jsonSchema"`
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

// MakeCollection convert to message
func (coll *MongoCollection) makeCollection() (collection *api.Collection) {
	collection = &api.Collection{
		Name:     coll.Name,
		Type:     coll.Type,
		ReadOnly: coll.Info.ReadOnly,
	}
	sch := coll.Options.Validator.JSONSchema
	collection.Schema = sch.makeSchema(coll.Name)
	return
}
