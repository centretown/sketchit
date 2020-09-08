package main

import "go.mongodb.org/mongo-driver/bson"

var delayCommand = bson.M{
	"title":       "Delay Command",
	"bsonType":    "object",
	"description": "delays for <duration>",
	"required":    []string{"duration"},
	"properties": bson.M{
		"duration": bson.M{
			"title":       "Duration",
			"bsonType":    "number",
			"description": "the time to do something (delay)",
		},
	},
}

var modeCommand = bson.M{
	"title":       "Mode Command",
	"bsonType":    "object",
	"description": "prepares <signal> pin <id> for <mode>",
	"required":    []string{"id", "signal", "mode"},
	"properties": bson.M{
		"id": bson.M{
			"title":       "Id",
			"bsonType":    "number",
			"description": "the pin Id number as defined by the device",
		},
		"signal": bson.M{
			"title":       "Signal",
			"bsonType":    "string",
			"description": "the Signal type analog/digital",
			"enum":        []string{"analog", "digital"},
		},
		"mode": bson.M{
			"title":       "Mode",
			"bsonType":    "string",
			"description": "the input/output Mode",
			"enum":        []string{"input", "output"},
		},
	},
}

var pinCommand = bson.M{
	"title":       "Pin Command",
	"bsonType":    "object",
	"description": "<inputs/outputs> from/to <analog/digital> pin <id>",
	"required":    []string{"id", "signal", "mode", "value"},
	"properties": bson.M{
		"id": bson.M{
			"title":       "Id",
			"bsonType":    "number",
			"description": "the pin number as defined on the device",
		},
		"signal": bson.M{
			"title":       "Signal",
			"bsonType":    "string",
			"description": "the type of signal analog/digital",
			"enum":        []string{"analog", "digital"},
		},
		"mode": bson.M{
			"title":       "Mode",
			"bsonType":    "string",
			"description": "the read/write mode",
			"enum":        []string{"input", "output"},
		},
		"value": bson.M{
			"title":       "Value output or input",
			"bsonType":    "number",
			"description": "the value to read or write",
		},
	},
}

var hallCommand = bson.M{
	"title":       "Hall Command",
	"bsonType":    "object",
	"description": "reads magnetic field to measurement",
	"required":    []string{"measurement"},
	"properties": bson.M{
		"measurement": bson.M{
			"title":       "Measurement",
			"bsonType":    "number",
			"description": "the measurement read",
		},
	},
}

var commandSchema = bson.M{
	"title":       "Command",
	"bsonType":    "object",
	"description": "the possible command parameters",
	"oneOf": bson.A{
		delayCommand,
		modeCommand,
		pinCommand,
		hallCommand,
	},
}
