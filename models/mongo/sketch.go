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
		"toolkit",
		"label",
		"device",
		"purpose",
		// "setup",
		// "loop",
	},
	// "uniqueItems": true,
	"properties": bson.M{
		"toolkit": bson.M{
			"title":       "Toolkit",
			"bsonType":    "string",
			"description": "the toolkit defines the set of operations available",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "the label assigned is unique for the toolkit",
		},
		"device": bson.M{
			"title":       "Device",
			"bsonType":    "string",
			"description": "the route to the device being operated",
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
	pathName:   "sketchIndex",
	pathUnique: true,
}

var sketchIndeces = []mongo.IndexModel{
	{
		Keys: bson.M{
			"toolkit": 1,
			"label":   1,
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
