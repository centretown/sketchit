PKG := "github.com/centretown/sketchit"
GOSOURCE := ${GOPATH}/src/$(PKG)
SERVER_OUT := $(GOSOURCE)"/bin/server"
CLIENT_OUT := $(GOSOURCE)"/bin/client"
BUILD_MONGO_SCHEMA_OUT := $(GOSOURCE)"/bin/build-mongo-schema"
API_OUT := $(GOSOURCE)"/api/sketchit.pb.go"
API_ACTION_OUT := $(GOSOURCE)"/api/action.pb.go"
GW_OUT := $(GOSOURCE)"/api/sketchit.pb.gw.go"
SWAG_OUT := $(GOSOURCE)"/api/sketchit.swagger.json"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
MONGO_SCHEMA_PKG_BUILD := "${PKG}/db-schema/mongo"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all api build_server build_client build_mongo_schema

all: build_server build_client build_mongo_schema

api/sketchit.pb.go: protos/sketchit.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketchit.proto 
	
	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketchit.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketchit.proto

api/action.pb.go: protos/action.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    action.proto 

	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src action.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api action.proto

	

test: ## run unit tests
	@client/curl_test.sh
	@go test ./...

api: api/sketchit.pb.go api/action.pb.go ## Auto-generate grpc go sources

dep: ## Get the dependencies
	@go get -v -d ./...

build_server: dep api ## Build the binary file for server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

build_client: dep api ## Build the binary file for client
	@go build -i -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

build_mongo_schema: ##
	@go build -i -v -o $(BUILD_MONGO_SCHEMA_OUT) $(MONGO_SCHEMA_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT) $(GW_OUT) $(SWAG_OUT) $(API_ACTION_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
