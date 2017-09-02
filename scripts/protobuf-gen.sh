#!/bin/bash
DIR=$GOPATH/src/github.com/ng-vu/go-grpc-sample
IMPORT="-I$DIR/protobuf \
    -I$DIR/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"

rm $DIR/pb/**/*.pb.go
rm $DIR/pb/**/*.pb.gw.go
rm $DIR/doc/*.swagger.json

for pkg in $(ls -c $DIR/protobuf | grep -v google); do
    PROTO=$DIR/protobuf/$pkg/*.proto
    protoc $IMPORT --go_out=plugins=grpc:$GOPATH/src/. $PROTO
    protoc $IMPORT --grpc-gateway_out=logtostderr=true:$DIR/pb $PROTO
    protoc $IMPORT --swagger_out=logtostderr=true:$DIR/protobuf $DIR/protobuf/$pkg/$pkg.proto

    mv $DIR/protobuf/$pkg/*.swagger.json $DIR/doc
done
