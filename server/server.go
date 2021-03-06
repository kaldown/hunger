package main

import (
	pb "../build/proto"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

const (
	address = "0.0.0.0"
)

var db *sql.DB

// 0 => PROD
// 1 => DEV
var MODE = 1

type server struct{}

// TODO: move
func getQuiz(db *sql.DB, quizTitle string) (string, error) {
	var title string
	err := db.QueryRow("SELECT id FROM quiz WHERE title = $1", quizTitle).Scan(&title)
	if err != nil {
		return "", err
	}
	return title, nil
}

// TODO: move
func setQuiz(db *sql.DB, quizTitle string) error {
	_, err := db.Exec(
		"INSERT INTO quiz (title) VALUES ($1)",
		quizTitle,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) GetQuiz(ctx context.Context, in *pb.QuizRequest) (*pb.QuizResponse, error) {
	quiz, err := getQuiz(db, in.Message)
	if err != nil {
		return nil, err
	}
	return &pb.QuizResponse{Message: quiz}, nil
}

func (s *server) SetQuiz(ctx context.Context, in *pb.QuizRequest) (*pb.QuizResponse, error) {
	err := setQuiz(db, in.Message)
	if err != nil {
		return &pb.QuizResponse{}, err
	}
	return &pb.QuizResponse{Message: "Created"}, nil
}

func main() {
	connStr := "user=postgres dbname=hunger-db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failsed to connect to db %v", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	sock := address + ":" + port
	log.Printf("Server is running on sock: %v", sock)

	lis, err := net.Listen("tcp", sock)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	switch MODE {
	case 0:
		log.Println("PROD mode")

		creds, err := credentials.NewServerTLSFromFile("../cert/server.crt", "../cert/server.key")
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		s := grpc.NewServer(grpc.Creds(creds))

		pb.RegisterQuizerServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	case 1:
		log.Println("DEV mode")

		s := grpc.NewServer()

		pb.RegisterQuizerServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

}
