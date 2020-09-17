package main

import "go.mongodb.org/mongo-driver/bson"

var actionSchema = bson.M{
	"title":       "Action",
	"bsonType":    "object",
	"description": "the action taken on the device",
	"required":    []string{"sequence", "type"},
	"uniqueItems": true,
	"properties": bson.M{
		"sequence": bson.M{
			"title":       "Sequence",
			"bsonType":    "number",
			"description": "the action sequence number",
		},
		"type": bson.M{
			"title":       "Type",
			"bsonType":    "string",
			"description": "the operation type",
		},
		"operation": operationSchema,
	},
}

var operationSchema = bson.M{
	"title":       "Operation",
	"bsonType":    "object",
	"description": "the operations available on this device",
	"oneOf": bson.A{
		delayOperation,
		modeOperation,
		pinOperation,
		hallOperation,
	},
}

var delayOperation = bson.M{
	"title":       "Delay",
	"bsonType":    "object",
	"description": "the delay operation pauses before the next action",
	"required":    []string{"duration"},
	"properties": bson.M{
		"duration": durationDef,
	},
}

var modeOperation = bson.M{
	"title":       "Mode",
	"bsonType":    "object",
	"description": "the mode operation prepares a pin for reading or writing",
	"required":    []string{"id", "signal", "mode"},
	"properties": bson.M{
		"id":     pinIDDef,
		"signal": signalDef,
		"mode":   modeDef,
	},
}

var pinOperation = bson.M{
	"title":       "Pin",
	"bsonType":    "object",
	"description": "the pin operation performs a read or write operation on the pin.",
	"required":    []string{"id", "signal", "mode", "value"},
	"properties": bson.M{
		"id":     pinIDDef,
		"signal": signalDef,
		"mode":   modeDef,
		"value":  valueDef,
	},
}

var hallOperation = bson.M{
	"title":       "Hall",
	"bsonType":    "object",
	"description": "reads magnetic field to measurement",
	"required":    []string{"measurement"},
	"properties": bson.M{
		"measurement": measurementDef,
	},
}

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
