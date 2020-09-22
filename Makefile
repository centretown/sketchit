PKG := "github.com/centretown/sketchit"
GOSOURCE := ${GOPATH}/src/$(PKG)
SERVER_OUT := $(GOSOURCE)"/bin/server"
CLIENT_OUT := $(GOSOURCE)"/bin/client"
BUILD_MONGO_SCHEMA_OUT := $(GOSOURCE)"/bin/models"
API_OUT := $(GOSOURCE)"/api/sketchit.pb.go"
API_ACTION_OUT := $(GOSOURCE)"/api/sketch.pb.go"
GW_OUT := $(GOSOURCE)"/api/sketchit.pb.gw.go"
SWAG_OUT := $(GOSOURCE)"/api/sketchit.swagger.json"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
MONGO_SCHEMA_PKG_BUILD := "${PKG}/models/mongo"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all api build_server build_client build_mongo_schema

all: build_server build_client build_mongo_schema
 
# generate sketchit services
api/sketchit.pb.go: protos/sketchit.proto
	# generate sketchit GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketchit.proto 
	
	# generate GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketchit.proto

	# generate Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketchit.proto

# generate sketch 
api/sketch.pb.go: protos/sketch.proto
	# generate sketch GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketch.proto 

	# generate GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketch.proto

	# generate Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketch.proto

# generate device
api/device.pb.go: protos/device.proto
	# generate device GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    device.proto 

	# generate GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src device.proto

	# generate Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api device.proto

# generate dictionary
api/dictionary.pb.go: protos/dictionary.proto
	# generate dictionary GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    dictionary.proto 

	# generate GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src dictionary.proto

	# generate Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api dictionary.proto

# generate deputy
api/deputy.pb.go: protos/deputy.proto
	# generate GRPC deputy services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    deputy.proto 

	# generate GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src deputy.proto

	# generate Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api deputy.proto

test: ## run unit tests
	@client/curl_test.sh
	@go test ./...

api: api/sketchit.pb.go api/sketch.pb.go api/device.pb.go api/deputy.pb.go api/dictionary.pb.go ## Auto-generate grpc go sources

dep: ## Get the dependencies
	@go get -v -d ./...

server: dep api ## Build the binary file for server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

client: dep api ## Build the binary file for client
	@go build -i -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

models: ##
	@go build -i -v -o $(BUILD_MONGO_SCHEMA_OUT) $(MONGO_SCHEMA_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT) $(GW_OUT) $(SWAG_OUT) $(API_ACTION_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
