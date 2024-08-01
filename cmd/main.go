package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sukhmai/fun/pkg/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	port = ":8080"
)

func main() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	sugar.Log(zapcore.InfoLevel, "test")

	r := mux.NewRouter()

	server, err := api.NewServer()
	if err != nil {
		sugar.Errorf("error creating server: %v", err)
		return
	}

	r.HandleFunc("/", server.HomeHandler)
	registerQuestionRoutes(r, server)
	registerProfileRoutes(r, server)

	handler := cors.AllowAll().Handler(r)

	http.Handle("/", handler)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sugar.Logf(zapcore.InfoLevel, "starting server on %s", port)
		err = http.ListenAndServe(port, nil)
		if err != nil {
			sugar.Errorf("error starting server: %v", err)
			return
		}
	}()

	// Block until a signal is received
	<-quit
	sugar.Log(zapcore.InfoLevel, "Shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = server.Close(ctxShutDown)
	if err != nil {
		sugar.Errorf("error closing server: %v", err)
	}
	sugar.Log(zapcore.InfoLevel, "Server gracefully stopped")
}

func registerQuestionRoutes(r *mux.Router, server *api.Server) {
	questions := r.PathPrefix("/questions").Subrouter()
	// questions.HandleFunc("/savereply/{userid:[0-9]+}/{questionnumber}", server.SaveReplyHandler)
	questions.HandleFunc("/saveanswer/{questionnumber}", server.SaveAnswerHandler)
	questions.HandleFunc("/getanswer/{questionnumber}", server.GetAnswerHandler)
	questions.HandleFunc("/getanswers", server.GetAnswersHandler)
	questions.HandleFunc("/getrandomquestion", server.GetRandomQuestionsHandler)
	questions.HandleFunc("/getquestions", server.GetQuestionsHandler)
}

func registerProfileRoutes(r *mux.Router, server *api.Server) {
	profile := r.PathPrefix("/profile").Subrouter()
	profile.HandleFunc("/uploadphoto", server.UploadProfilePhotoHandler)
	profile.HandleFunc("/getphoto", server.GetProfilePhotoHandler)
	profile.HandleFunc("/saveuser", server.SaveUserHandler)
	profile.HandleFunc("/defaultphoto", server.GetDefaultPhotoHandler)
}
