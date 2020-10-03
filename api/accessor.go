package api

import (
	"fmt"

	"github.com/centretown/sketchit/info"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// Accessor -
type Accessor struct {
	MakeItem func() protoreflect.ProtoMessage
	Parent   string
}

// ProtoAccessor -
type protoAccessor struct {
	makeItem   func() protoreflect.ProtoMessage
	parentName string
}

const (
	devicesName  = "devices"
	sketchesName = "sketches"
	deputiesName = "deputies"
	sectorsName  = "sectors"
	toolkitsName = "toolkits"
)

var messages = make(map[string]*protoAccessor)

func init() {
	messages[devicesName] = &protoAccessor{
		makeItem:   func() protoreflect.ProtoMessage { return &Device{} },
		parentName: sectorsName,
	}
	messages[sketchesName] = &protoAccessor{
		makeItem:   func() protoreflect.ProtoMessage { return &Sketch{} },
		parentName: toolkitsName,
	}
}

// NewAccessor -
func NewAccessor(route ...string) (accessor *Accessor, err error) {
	if len(route) < 1 {
		err = info.Inform(err, ErrNoRoute, route)
		return
	}
	name := route[0]
	proto, ok := messages[name]
	if !ok {
		err = info.Inform(err, ErrNoAccessor, route)
		return
	}
	accessor = &Accessor{
		MakeItem: proto.makeItem,
		Parent:   buildParent(proto.parentName, name, route...),
	}
	return
}

func buildParent(parentName,
	collectionName string, route ...string) (parent string) {

	switch len(route) {
	case 0:
	case 1:
		parent = fmt.Sprintf("%s/*/%s", parentName, collectionName)
	case 2:
		parent = fmt.Sprintf("%s/%s/%s",
			parentName, route[1], collectionName)
	default:
		parent = fmt.Sprintf("%s/%s/%s/%s",
			parentName, route[1], collectionName, route[2])
	}

	return
}
