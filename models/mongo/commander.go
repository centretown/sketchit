package main

import "go.mongodb.org/mongo-driver/bson"

var commanderSchema = bson.M{
	"title":       "Commander",
	"bsonType":    "object",
	"description": "the commander api",
	"required":    []string{"label", "version", "commands"},
	"uniqueItems": true,
	"properties": bson.M{
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "a unique label identifying the commander release",
		},
		"version": bson.M{
			"title":       "Version",
			"bsonType":    "string",
			"description": "the version is unique to each label and changes when a command is added, changed or removed.",
		},
		"commands": bson.M{
			"title":       "Commands",
			"bsonType":    "array",
			"description": "the list of available commands ordered by their verb",
			"items":       commandSchema,
		},
		// build at runtime
		"aliases": bson.M{
			"title":       "Aliases",
			"bsonType":    "array",
			"description": "the map of commands keyed on the verbs and their alternates.",
			"items":       aliasSchema,
		},
		// built at runtime
		"dictionary": bson.M{
			"title":       "Dictionary",
			"bsonType":    "array",
			"description": "the dictionary provides a pathway to the vocabulary of models derived from the api.",
			"items":       dictionarySchema,
		},
		// built at runtime
		"gallery": bson.M{
			"title":       "Gallery",
			"bsonType":    "array",
			"description": "the gallery provides",
			"items":       gallerySchema,
		},
	},
}
