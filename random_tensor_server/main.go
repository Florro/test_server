package main

import (
	"fmt"
	tf "github.com/florro/test_server/proto/tensorflow"
	pb "github.com/florro/test_server/proto/tensorflow_serving"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Predict(ctx context.Context, req *pb.PredictRequest) (*pb.PredictResponse, error) {

	var size int64 = 256

	var tensor *tf.TensorProto = createrandomtensor(size)

	tmpmap := make(map[string]*tf.TensorProto)
	tmpmap["bools"] = tensor
	var resp = &pb.PredictResponse{Outputs: tmpmap}

	return resp, nil

}

func createrandomtensor(size int64) *tf.TensorProto {

	Shape := &tf.TensorShapeProto{
		Dim: []*tf.TensorShapeProto_Dim{
			&tf.TensorShapeProto_Dim{Name: "x", Size: size},
		},
	}
	fmt.Println(Shape)

	var tmp = make([]bool, size)
	for i := 0; i < len(tmp); i++ {
		if rand.Intn(2) != 0 {
			tmp[i] = true
		}
	}

	tensor := &tf.TensorProto{
		BoolVal:     tmp,
		TensorShape: Shape,
		Dtype:       tf.DataType_DT_BOOL,
	}

	return tensor
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPredictionServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)

	}
}
