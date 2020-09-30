package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Collector -
type Collector struct {
	Collection *mongo.Collection
	NewItem    func() protoreflect.ProtoMessage
	Filter     func(parent string) bson.D
}

// CollectorKeys -
const (
	ParentName int = iota
	ParentLabel
	CollectionName
	CollectionLabel
)
