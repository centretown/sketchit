package main

import (
	"github.com/centretown/sketchit/api"
	"go.mongodb.org/mongo-driver/bson"
)

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
			"title": "Flag",
			"enum": bson.A{
				api.Feature_f,
				api.Feature_d,
				api.Feature_auto,
			},
			"description": "",
		},
		"label": bson.M{
			"title":       "Label",
			"bsonType":    "string",
			"description": "unique label to reference feature",
		},
		"summary": summaryDef,
	},
}
