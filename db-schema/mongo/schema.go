package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var actionSchema = bson.M{
	"title":       "Action",
	"bsonType":    "object",
	"description": "the action to take",
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

var processTrait = &indexTrait{
	pathUnique: true,
	pathName:   "pathIndex",
}

var processIndeces = []mongo.IndexModel{
	{
		Keys: bson.M{
			"model": 1,
			"label": 1,
		},
		Options: &options.IndexOptions{
			Name:   &processTrait.pathName,
			Unique: &processTrait.pathUnique,
		},
	},
	{
		Keys: bson.M{
			"deviceKey": 1,
		},
	},
}

var processSchema = bson.M{
	"bsonType": "object",
	"required": []string{
		"model",
		"label",
		"deviceKey",
		"purpose",
	},
	"properties": bson.M{
		"model": bson.M{
			"title":       "Model",
			"bsonType":    "string",
			"description": "the device model type",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "the name of the process",
		},
		"deviceKey": bson.M{
			"title":       "Device Key",
			"bsonType":    "string",
			"description": "the key for the device",
		},
		"purpose": bson.M{
			"title":       "Purpose",
			"bsonType":    "string",
			"description": "the purpose or intent",
		},
		"setup": bson.M{
			"title":       "Setup",
			"bsonType":    "array",
			"description": "the actions taken during the setup stage",
			"items":       actionSchema,
		},
		"loop": bson.M{
			"title":       "Loop",
			"bsonType":    "array",
			"description": "the actions taken during the processing loop",
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
			"domain": 1,
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
		"description": "the pin definition",
		"required":    []string{"id", "label", "purpose"},
		"uniqueItems": true,
		"properties": bson.M{
			"id": bson.M{
				"title":       "Id",
				"bsonType":    "number",
				"description": "the pin number",
			},
			"label": bson.M{
				"title":       "Label",
				"bsonType":    "string",
				"description": "the label assigned to the pin",
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
	"bsonType": "object",
	"required": []string{"domain", "label", "model"},
	"properties": bson.M{
		"domain": bson.M{
			"title":       "Domain",
			"bsonType":    "string",
			"description": "the name of the domain",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "the name of the device",
		},
		"model": bson.M{
			"title":       "Model",
			"bsonType":    "string",
			"description": "the device model type",
		},
		"pins": pinSchema,
	},
}
