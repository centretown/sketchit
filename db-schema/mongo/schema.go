package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var actionSchema = bson.M{
	"title":       "Action",
	"bsonType":    "object",
	"description": "the action taken on the device",
	"required":    []string{"sequence", "type"},
	"uniqueItems": true,
	"properties": bson.M{
		"sequence": bson.M{
			"title":       "Sequence",
			"bsonType":    "number",
			"description": "the action sequence number",
		},
		"type": bson.M{
			"title":       "Type",
			"bsonType":    "string",
			"description": "the command type",
		},
		"command": commandSchema,
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
			"deviceKey": 1,
		},
	},
}

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

var deviceSchema = bson.M{
	"title":    "Devices",
	"bsonType": "object",
	"required": []string{"sector", "label", "model"},
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
