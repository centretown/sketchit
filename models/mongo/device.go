package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var deviceSchema = bson.M{
	"title":       "Devices",
	"bsonType":    "object",
	"description": "controller devices",
	"required":    []string{"sector", "label", "model"},
	"properties": bson.M{
		"sector": bson.M{
			"title":       "Sector",
			"bsonType":    "string",
			"description": "the sector groups interconnected devices",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "the label assigned is unique in the sector",
		},
		"model": bson.M{
			"title":       "Model",
			"bsonType":    "string",
			"description": "the device model type",
		},
		"pins": pinSchema,
	},
}

var deviceTrait = &indexTrait{
	pathUnique: true,
	pathName:   "pathIndex",
}

var deviceIndeces = []mongo.IndexModel{
	{
		Keys: bson.M{
			"sector": 1,
			"label":  1,
		},
		Options: &options.IndexOptions{
			Name:   &deviceTrait.pathName,
			Unique: &deviceTrait.pathUnique,
		},
	},
	{
		Keys: bson.M{
			"model": 1,
		},
	},
}

var pinSchema = bson.M{
	"title":       "Pins",
	"bsonType":    "array",
	"description": "the active pins on this device",
	"items": bson.M{
		"title":       "Pin",
		"bsonType":    "object",
		"description": "the pin as defined by its purpose",
		"required":    []string{"id", "label", "purpose"},
		"uniqueItems": true,
		"properties": bson.M{
			"id": pinIDDef,
			"label": bson.M{
				"title":       "Label",
				"bsonType":    "string",
				"description": "the label assigned is unique for this device",
			},
			"purpose": bson.M{
				"title":       "Purpose",
				"bsonType":    "string",
				"description": "the purpose for this pin",
			},
		},
	},
}
