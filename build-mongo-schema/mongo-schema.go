package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var commandSchema = bson.M{
	"title":       "Command",
	"bsonType":    "object",
	"description": "the possible command parameters",
	"properties": bson.M{
		"id": bson.M{
			"bsonType":    "number",
			"description": "the pin number",
		},
		"signal": bson.M{
			"bsonType":    "string",
			"description": "the signal type analog/digital, string",
		},
		"mode": bson.M{
			"bsonType":    "string",
			"description": "the i/o type input/output, string",
		},
		"value": bson.M{
			"bsonType":    "number",
			"description": "the value written or read",
		},
		"measurement": bson.M{
			"bsonType":    "number",
			"description": "the measurement read",
		},
		"duration": bson.M{
			"bsonType":    "number",
			"description": "the time to do something (delay)",
		},
	},
}

var actionSchema = bson.M{
	"title":       "Action",
	"bsonType":    "object",
	"description": "the action definition",
	"required":    []string{"sequence", "type"},
	"uniqueItems": true,
	"properties": bson.M{
		"sequence": bson.M{
			"bsonType":    "number",
			"description": "the sequence operation, required number",
		},
		"type": bson.M{
			"bsonType":    "string",
			"description": "the command type, required string",
		},
		"command": commandSchema,
	},
}

type indexTrait struct {
	idUnique   bool
	pathUnique bool
	pathName   string
}

var processTrait = &indexTrait{
	idUnique:   true,
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
			Name:   &deviceTrait.pathName,
			Unique: &deviceTrait.pathUnique,
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
	},
	"properties": bson.M{
		"model": bson.M{
			"bsonType":    "string",
			"description": "the device model type, required string",
		},
		"label": bson.M{
			"bsonType":    "string",
			"description": "the name of the process, required string",
		},
		"deviceKey": bson.M{
			"bsonType":    "string",
			"description": "the key for the device, required string",
		},
		"purpose": bson.M{
			"bsonType":    "string",
			"description": "the purpose or intent, string",
		},
		"setup": bson.M{
			"bsonType":    "array",
			"description": "the actions taken during the setup stage",
			"items":       actionSchema,
		},
		"loop": bson.M{
			"bsonType":    "array",
			"description": "the actions taken during the processing loop",
			"items":       actionSchema,
		},
	},
}

var deviceTrait = &indexTrait{
	idUnique:   true,
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

var deviceSchema = bson.M{
	"bsonType": "object",
	"required": []string{"domain", "label", "model"},
	"properties": bson.M{
		"domain": bson.M{
			"bsonType":    "string",
			"description": "the name of the domain, required string",
		},
		"label": bson.M{
			"bsonType":    "string",
			"description": "the name of the device, required string",
		},
		"model": bson.M{
			"bsonType":    "string",
			"description": "the device model type, required string",
		},
		"pins": bson.M{
			"bsonType":    "array",
			"description": "the managed pins on this device",
			"items": bson.M{
				"title":       "Pin",
				"bsonType":    "object",
				"description": "the pin definition",
				"required":    []string{"id", "label"},
				"uniqueItems": true,
				"properties": bson.M{
					"id": bson.M{
						"bsonType":    "number",
						"description": "the pin number, required number",
					},
					"label": bson.M{
						"bsonType":    "string",
						"description": "the label assigned to the pin, required string",
					},
					"purpose": bson.M{
						"bsonType":    "string",
						"description": "the purpose or intent, string",
					},
				},
			},
		},
	},
}
