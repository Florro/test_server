#protoc -I ./ --go_out=. predict.proto
protoc --go_out=plugins=grpc,import_path=github.com/florro/test_server/proto/helloworld:. helloworld.proto
