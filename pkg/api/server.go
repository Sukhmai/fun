package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/sukhmai/fun/pkg/db"
	"github.com/sukhmai/fun/pkg/redis"
	"go.uber.org/zap"
)

type Server struct {
	dbClient    *db.DBClient
	redisClient *redis.RedisClient
	logger      *zap.SugaredLogger
}

func NewServer() (*Server, error) {
	prod, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	password := os.Getenv("DB_PASSWORD")
	mode := os.Getenv("MODE")
	if password == "" {
		return nil, errors.New("DB_PASSWORD environment variable not set")
	}
	var connString string
	if mode == "production" {
		connString = "postgres://fun:" + password + "@localhost:5432/sample"
	} else {
		connString = "postgres://postgres:" + password + "@localhost:5432/postgres"
	}
	dbClient, err := db.NewClient(connString)
	if err != nil {
		return nil, err
	}
	logger := prod.Sugar()
	return &Server{
		dbClient:    dbClient,
		redisClient: redis.NewClient(),
		logger:      logger,
	}, nil

}

func (s *Server) Close(ctx context.Context) error {
	err := s.redisClient.Close()
	if err != nil {
		return err
	}
	s.dbClient.Close()
	if err != nil {
		return err
	}
	return nil
}

type Response struct {
	Message string `json:"message"`
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func (s *Server) returnError(w http.ResponseWriter, message string, statusCode int) error {
	s.logger.Errorf(message)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Message: message})
	return errors.New(message)
}
