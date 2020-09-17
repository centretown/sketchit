package main

import "go.mongodb.org/mongo-driver/bson"

var commandSchema = bson.M{
	"title":       "Command",
	"bsonType":    "object",
	"description": "the action taken on the device",
	"required":    []string{"sequence", "type"},
	"uniqueItems": true,
	"properties": bson.M{
		"verb":         verbDef,
		"summary":      summaryDef,
		"arguments":    argumentDef,
		"presentation": presentationDef,
	},
}

var argumentDef = bson.M{
	"title":       "Arguments",
	"bsonType":    "array",
	"description": "list input parameters",
	"items": bson.M{
		"title":       "Argument",
		"bsonType":    "string",
		"description": "parsed argument",
	},
}

var summaryDef = bson.M{
	"usage": bson.M{
		"title":       "Usage",
		"bsonType":    "string",
		"description": "basic operation",
	},
	"syntax": bson.M{
		"title":       "Syntax",
		"bsonType":    "string",
		"description": "rules for use",
	},
	"examples": examplesSchema,
}

var examplesSchema = bson.M{
	"title":       "Examples",
	"bsonType":    "array",
	"description": "list of typical cases",
	"items": bson.M{
		"title":       "Example",
		"bsonType":    "string",
		"description": "use case",
	},
}

var verbDef = bson.M{
	"title":       "Verb",
	"bsonType":    "string",
	"description": "available actions to take ",
	"enum": []string{
		"help",
		"list",
		"goto",
		"save",
		"remove",
		"format",
		"detail",
		"auto",
		"hello",
		"exit",
	},
}
