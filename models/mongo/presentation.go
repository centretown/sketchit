package main

import "go.mongodb.org/mongo-driver/bson"

var presentationDef = bson.M{
	"title":       "Presentation",
	"bsonType":    "object",
	"description": "the presentation features selected",
	"required":    []string{"format", "projection", "auto"},
	"properties": bson.M{
		"format":     formatDef,
		"projection": projectionDef,
		"auto":       autoDef,
	},
}

var formatDef = bson.M{
	"title":       "Format",
	"bsonType":    "object",
	"description": "the output formats available",
	"oneOf": bson.A{
		"yaml",
		"json",
		"xml",
	},
}

var projectionDef = bson.M{
	"title":       "Projection",
	"bsonType":    "array",
	"description": "the levels of detail",
	"items": bson.M{
		"title":       "Argument",
		"bsonType":    "string",
		"description": "parsed argument",
		"oneOf": bson.A{
			"full",
			"summary",
			"brief",
		},
	},
}

var autoDef = bson.M{
	"title":       "Auto",
	"bsonType":    "object",
	"description": "auto reply feature",
	"oneOf": bson.A{
		"off",
		"y",
		"n",
	},
}
