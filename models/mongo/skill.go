package main

import "go.mongodb.org/mongo-driver/bson"

var skillDef = bson.M{
	"title":       "Skill",
	"bsonType":    "object",
	"description": "The ability to perform a task.",
	"required": []string{
		"task",
	},
	"uniqueItems": true,
	"properties": bson.M{
		"task":       taskDef,
		"alternates": alternatesDef,
		"summary":    summaryDef,
	},
}

var taskDef = bson.M{
	"title":       "Task",
	"bsonType":    "string",
	"description": "A job item with an outcome. A task is identifed by a command issued to the deputy, which finds a matching skill to perform the task. The outcome is relayed to the issuer.",
	"enum": []string{
		"hello",
		"help",
		"skills",
		"features",
		"list",
		"goto",
		"edit",
		"save",
		"remove",
		"exit",
	},
}

var alternatesDef = bson.M{
	"title":       "Alternates",
	"bsonType":    "array",
	"description": "list of synonyms",
	"items": bson.M{
		"title":       "Alternate",
		"bsonType":    "string",
		"description": "synonym",
	},
}

var summaryDef = bson.M{
	"title":       "Summary",
	"bsonType":    "object",
	"description": "summary description",
	"properties": bson.M{
		"usage": bson.M{
			"title":       "Usage",
			"bsonType":    "string",
			"description": "usage statement",
		},
		"syntax": bson.M{
			"title":       "Syntax",
			"bsonType":    "string",
			"description": "usage rules",
		},
		"examples": examplesDef,
	},
}

var examplesDef = bson.M{
	"title":       "Examples",
	"bsonType":    "array",
	"description": "use cases",
	"items": bson.M{
		"title":       "Example",
		"bsonType":    "string",
		"description": "use case",
	},
}
