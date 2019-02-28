package main

import (
	pb "../build/generated"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	address = "0.0.0.0"
)

var quizes = map[string]string{}

type server struct{}

func (s *server) GetQuiz(ctx context.Context, in *pb.QuizRequest) (*pb.QuizResponse, error) {
	log.Printf("getQuiz: %v", in)
	quiz, ok := quizes[in.Message]
	if !ok {
		return nil, errors.New("No such quiz")
	}
	return &pb.QuizResponse{Message: quiz}, nil
}

func (s *server) SetQuiz(ctx context.Context, in *pb.QuizRequest) (*pb.QuizResponse, error) {
	log.Printf("setQuiz: %v", in)
	quizes[in.Message] = in.Message
	return &pb.QuizResponse{Message: "Quiz was saved"}, nil
}

func main() {
	port := os.Getenv("PORT")
	log.Printf("Get $PORT: %v", port)
	sock := address + ":" + port
	log.Printf("Server is running on sock: %v", sock)
	lis, err := net.Listen("tcp", sock)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQuizerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
