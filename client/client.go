package main

import (
	pb "../build/generated"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:5001"
	defaultName = "salo"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewQuizerClient(conn)

	quiz := defaultName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetQuiz(ctx, &pb.QuizRequest{Message: quiz})
	if err != nil {
		log.Printf("GetQuiz Request failed: %v", err)

		// TEST: try to fillup and ask again
		r2, err2 := c.SetQuiz(ctx, &pb.QuizRequest{Message: quiz})
		if err2 != nil {
			log.Fatalf("SetQuiz Request failed: %v", err2)
		}
		log.Printf("Response: %v", r2.Message)

		// ask again
		r3, err3 := c.GetQuiz(ctx, &pb.QuizRequest{Message: quiz})
		if err3 != nil {
			log.Fatalf("GetQuiz Second Request failed: %v", err3)
		}
		log.Printf("Response: %v", r3.Message)
		return
		// END Test

	}
	log.Printf("Response: %v", r.Message)
}
