package main

import "go.mongodb.org/mongo-driver/bson"

var featureDef = bson.M{
	"title":       "Feature",
	"bsonType":    "object",
	"description": "the deputy's feature set",
	"required": []string{
		"flag",
		"label",
	},
	"properties": bson.M{
		"flag": bson.M{
			"title":       "Flag",
			"bsonType":    "string",
			"description": "",
			"enum": bson.A{
				"f",
				"d",
				"auto",
			},
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "unique label to reference feature",
		},
		"summary": summaryDef,
	},
}
