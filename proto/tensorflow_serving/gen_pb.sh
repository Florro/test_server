#protoc -I ./ --go_out=. predict.proto
protoc --proto_path=. -I$GOPATH/src -I$GOPATH/src/github.com/florro/test_server/proto/tensorflow --go_out=plugins=grpc:. *.proto
