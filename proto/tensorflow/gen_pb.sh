#protoc -I ./ --go_out=. predict.proto
protoc --proto_path=. --proto_path=/home/florian/go/src --go_out=plugins=grpc:. *.proto
