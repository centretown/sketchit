# api/sketch.pb.go: protos/sketch.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    sketch.proto 

	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src sketch.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api sketch.proto

# api/device.pb.go: protos/device.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    device.proto 

	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src device.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api device.proto

# api/collection.pb.go: protos/collection.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    collection.proto 

	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src collection.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api collection.proto

# api/deputy.pb.go: protos/deputy.proto
	@protoc -I $(GOSOURCE)/protos \
        --go_out=plugins=grpc:$(GOPATH)/src\
    deputy.proto 

	# create reverse proxy
	@protoc -I protos --grpc-gateway_out=logtostderr=true:${GOPATH}/src deputy.proto

	# create Open API doc for REST interface
	@protoc -I protos --swagger_out=logtostderr=true:$(GOSOURCE)/api deputy.proto

