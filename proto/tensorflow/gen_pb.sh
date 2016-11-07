#protoc -I ./ --go_out=. predict.proto
protoc --proto_path=. --proto_path=$GOPATH/src --go_out=plugins=grpc:. *.proto
