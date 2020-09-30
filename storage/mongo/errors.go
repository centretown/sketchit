package storage

import "errors"

// api support
var (
	ErrMongoInsert     = errors.New("failed to insert")
	ErrBsonMarshall    = errors.New("failed to create json")
	ErrCollection      = errors.New("failed to get collection")
	ErrFind            = errors.New("failed to find")
	ErrCreate          = errors.New("failed to create")
	ErrDecode          = errors.New("failed to decode")
	ErrDecodeAny       = errors.New("failed to decode any")
	ErrNoMatch         = errors.New("No records matched")
	ErrNotEnough       = errors.New("Not enough arguments")
	ErrCollectionNames = errors.New("failed to get collection names")
)

// connection support
var (
	ErrMongoClient  = errors.New("failed to create client")
	ErrMongoConnect = errors.New("failed to connect")
	ErrDisconnect   = errors.New("failed to disconnect")
)
