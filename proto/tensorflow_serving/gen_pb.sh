#protoc -I ./ --go_out=. predict.proto
protoc --proto_path=. -I/home/florian/go/src -I/home/florian/go/src/github.com/florro/test_server/proto/tensorflow --go_out=plugins=grpc:. *.proto
