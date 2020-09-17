package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sketchSchema = bson.M{
	"title":       "Sketches",
	"bsonType":    "object",
	"description": "the sketch defines the actions taken on a device",
	"required": []string{
		"model",
		"label",
		"device",
		"purpose",
		"setup",
		"loop",
	},
	"properties": bson.M{
		"model": bson.M{
			"title":       "Model",
			"bsonType":    "string",
			"description": "the model defines the scope of actions that can be taken",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "the label assigned is unique for the model",
		},
		"device": bson.M{
			"title":       "Device",
			"bsonType":    "string",
			"description": "the device acted on by this sketch",
		},
		"purpose": bson.M{
			"title":       "Purpose",
			"bsonType":    "string",
			"description": "the purpose of this sketch",
		},
		"setup": bson.M{
			"title":       "Setup",
			"bsonType":    "array",
			"description": "the setup stage list of actions",
			"items":       actionSchema,
		},
		"loop": bson.M{
			"title":       "Loop",
			"bsonType":    "array",
			"description": "the loop cycle list of actions",
			"items":       actionSchema,
		},
	},
}

type indexTrait struct {
	pathUnique bool
	pathName   string
}

var sketchTrait = &indexTrait{
	pathUnique: true,
	pathName:   "pathIndex",
}

var sketchIndeces = []mongo.IndexModel{
	{
		Keys: bson.M{
			"model": 1,
			"label": 1,
		},
		Options: &options.IndexOptions{
			Name:   &sketchTrait.pathName,
			Unique: &sketchTrait.pathUnique,
		},
	},
	{
		Keys: bson.M{
			"device": 1,
		},
	},
}
