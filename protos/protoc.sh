# go stuff
GOSOURCE=$GOPATH/src/github.com/centretown/sketchit/foreign
# for now this is quick so compile the languages i prefer to use
protoc --cpp_out=../foreign/cpp\
        --python_out=../foreign/python\
        --js_out=import_style=commonjs:../foreign/js \
        --grpc-web_out=import_style=commonjs,mode=grpcwebtext:../foreign/js\
        --js_out=import_style=typescript:../foreign/ts \
        --grpc-web_out=import_style=typescript,mode=grpcwebtext:../foreign/ts\
    device.proto 

echo $GOSOURCE
ls $GOSOURCE
