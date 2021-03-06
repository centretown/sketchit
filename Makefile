PKG := "github.com/centretown/sketchit"
GOSOURCE := ${GOPATH}/src/$(PKG)
SERVER_OUT := $(GOSOURCE)"/bin/server"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_OUT := $(GOSOURCE)"/bin/client"
CLIENT_PKG_BUILD := "${PKG}/client"
MONGO_SCHEMA_OUT := $(GOSOURCE)"/bin/models"
MONGO_SCHEMA_PKG_BUILD := "${PKG}/models/mongo"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

all: api server client models
.PHONY: all 
 
# generate sketchit services
api/sketchit.pb.go: protos/sketchit.proto
	# sketchit GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketchit.proto 
	
	# sketchit GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketchit.proto

	# sketchit Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketchit.proto

# generate sketch 
api/sketch.pb.go: protos/sketch.proto
	# sketch GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketch.proto 

	# sketch GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketch.proto

	# sketch Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketch.proto

# generate device
api/device.pb.go: protos/device.proto
	# device GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    device.proto 

	# device GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src device.proto

	# device Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api device.proto

# generate collection
api/collection.pb.go: protos/collection.proto
	# collection GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    collection.proto 

	# collection GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src collection.proto

	# collection Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api collection.proto

# generate deputy
api/deputy.pb.go: protos/deputy.proto
	# deputy GRPC services and protocol buffers
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    deputy.proto 

	# deputy GRPC Gateway reverse proxy for REST interface
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src deputy.proto

	# deputy Swagger Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api deputy.proto

test: ## Run unit tests
	@client/curl_test.sh
	@go test ./...

api: api/sketchit.pb.go api/sketch.pb.go api/device.pb.go api/deputy.pb.go api/collection.pb.go ## Auto-generate grpc go sources
# api: api/*.go ## Auto-generate grpc go sources

dep: ## Get the dependencies
	@go get -v -d ./...

$(MONGO_SCHEMA_OUT):
	# build models
	@go build -i -v -o $(MONGO_SCHEMA_OUT) $(MONGO_SCHEMA_PKG_BUILD)

$(SERVER_OUT): 
	# build server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

$(CLIENT_OUT):
	# build client
	@go build -i -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

server: $(SERVER_OUT) ## Build the server binary

client: $(CLIENT_OUT) ## Build the client binary

models: $(MONGO_SCHEMA_OUT) ## Build models binary

clean: ## Remove previous builds
	# remove binaries
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(MONGO_SCHEMA_OUT)
	# remove generated
	@rm api/sketchit.pb.go api/sketch.pb.go api/device.pb.go api/deputy.pb.go api/collection.pb.go

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
