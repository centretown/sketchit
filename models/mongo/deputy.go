package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var deputySchema = bson.M{
	"title":       "Deputy",
	"bsonType":    "object",
	"description": "The deputy performs tasks. Tasks are job items with an outcome. A task is identifed by a command issued to the deputy. The task, location and presentation features are identified. A matching skill is found perform the task. The outcome is relayed to the issuer.",
	// "required": []string{
	// 	"label",
	// 	"version",
	// 	"skills",
	// 	"features",
	// },
	// "uniqueItems": true,
	"properties": bson.M{
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "A unique label identifying the release. A release is created when skills are added, changed or removed.",
		},
		"version": bson.M{
			"title":       "Version",
			"bsonType":    "string",
			"description": "The version is unique to each label and is incremented when non breaking modifications are made.",
		},
		"skills": bson.M{
			"title":       "Skills",
			"bsonType":    "array",
			"description": "The available skills ordered by task.",
			"items":       skillDef,
		},
		"features": bson.M{
			"title":       "Features",
			"bsonType":    "array",
			"description": "The available features ordered by their label.",
			"items":       featureDef,
		},

		// // build at runtime
		// "aliases": bson.M{
		// 	"title":       "Aliases",
		// 	"bsonType":    "array",
		// 	"description": "the map of skills keyed on the verbs and their alternates.",
		// 	"items":       aliasSchema,
		// },
		// // built at runtime
		// "dictionary": bson.M{
		// 	"title":       "Dictionary",
		// 	"bsonType":    "array",
		// 	"description": "the dictionary provides a pathway to the vocabulary of models derived from the api.",
		// 	"items":       dictionarySchema,
		// },
		// // built at runtime
		// "gallery": bson.M{
		// 	"title":       "Gallery",
		// 	"bsonType":    "array",
		// 	"description": "the gallery provides",
		// 	"items":       gallerySchema,
		// },
	},
}

var deputyTrait = &indexTrait{
	pathUnique: true,
	pathName:   "pathIndex",
}

// version in descending order, FindOne should retrieve latest version.
var deputyIndeces = []mongo.IndexModel{
	{
		Keys: bson.M{
			"label":   1,
			"version": -1,
		},
		Options: &options.IndexOptions{
			Name:   &deputyTrait.pathName,
			Unique: &deputyTrait.pathUnique,
		},
	},
}
