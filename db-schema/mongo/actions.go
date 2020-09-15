package main

import "go.mongodb.org/mongo-driver/bson"

var signalDef = bson.M{
	"title":       "Signal",
	"bsonType":    "string",
	"description": "the signal type as defined by the device",
	"enum":        []string{"analog", "digital"},
}

var modeDef = bson.M{
	"title":       "Mode",
	"bsonType":    "string",
	"description": "the mode of operation",
	"enum":        []string{"input", "output"},
}

var valueDef = bson.M{
	"title":       "Value",
	"bsonType":    "number",
	"description": "the value to read or write",
}

var durationDef = bson.M{
	"title":       "Duration",
	"bsonType":    "number",
	"description": "the duration in milli-seconds",
}

var pinIDDef = bson.M{
	"title":       "Id",
	"bsonType":    "number",
	"description": "the id number of the pin as defined by this device",
}
var measurementDef = bson.M{
	"title":       "Measurement",
	"bsonType":    "number",
	"description": "the measurement read",
}

var delayCommand = bson.M{
	"title":       "Delay",
	"bsonType":    "object",
	"description": "the delay command pauses before the next action",
	"required":    []string{"duration"},
	"properties": bson.M{
		"duration": durationDef,
	},
}

var modeCommand = bson.M{
	"title":       "Mode",
	"bsonType":    "object",
	"description": "the mode command prepares a pin for reading or writing",
	"required":    []string{"id", "signal", "mode"},
	"properties": bson.M{
		"id":     pinIDDef,
		"signal": signalDef,
		"mode":   modeDef,
	},
}

var pinCommand = bson.M{
	"title":       "Pin",
	"bsonType":    "object",
	"description": "the pin command performs a read or write operation on the pin.",
	"required":    []string{"id", "signal", "mode", "value"},
	"properties": bson.M{
		"id":     pinIDDef,
		"signal": signalDef,
		"mode":   modeDef,
		"value":  valueDef,
	},
}

var hallCommand = bson.M{
	"title":       "Hall",
	"bsonType":    "object",
	"description": "reads magnetic field to measurement",
	"required":    []string{"measurement"},
	"properties": bson.M{
		"measurement": measurementDef,
	},
}

var commandSchema = bson.M{
	"title":       "Command",
	"bsonType":    "object",
	"description": "the command options available to this device",
	"oneOf": bson.A{
		delayCommand,
		modeCommand,
		pinCommand,
		hallCommand,
	},
}
