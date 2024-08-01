package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sukhmai/fun/pkg/questions"
)

type AnswerRequest struct {
	Answer string `json:"answer"`
}

type ReplyRequest struct {
	Reply string `json:"reply"`
}

type AnswersResponse struct {
	Answers map[string]string `json:"answers"`
}

func (s *Server) SaveAnswerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionNumString := vars["questionnumber"]
	userId, err := s.profileHandler(w, r)
	if err != nil {
		return
	}
	questionNum, err := strconv.Atoi(questionNumString)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error converting question number to int: %v", err), http.StatusBadRequest)
		return
	}
	var data AnswerRequest

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.returnError(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	if err = s.dbClient.SaveAnswer(context.Background(), userId, data.Answer, questionNum); err != nil {
		s.returnError(w, fmt.Sprintf("error saving answer to db: %v", err), http.StatusInternalServerError)
		return

	}
	response := Response{Message: "Answer saved!"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode((response))

}

func (s *Server) GetAnswerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := s.profileHandler(w, r)
	if err != nil {
		return
	}
	questionNumString := vars["questionnumber"]
	questionNum, err := strconv.Atoi(questionNumString)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error converting question number to int: %v", err), http.StatusBadRequest)
		return
	}
	answer, err := s.dbClient.GetAnswer(context.Background(), userId, questionNum)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error getting answer from db: %v", err), http.StatusInternalServerError)
		return
	}
	response := Response{Message: answer}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) GetAnswersHandler(w http.ResponseWriter, r *http.Request) {
	answers, err := s.dbClient.GetAllAnswers(context.Background())
	if err != nil {
		s.returnError(w, fmt.Sprintf("error getting answers from db: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(answers)
}

func (s *Server) SaveReplyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId1, err := s.profileHandler(w, r)
	if err != nil {
		return
	}
	userId2 := vars["userid"]
	questionNumString := vars["questionnumber"]
	questionNum, err := strconv.Atoi(questionNumString)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error converting question number to int: %v", err), http.StatusBadRequest)
		return
	}
	replyNumString := vars["replynumber"]

	replyNum, err := strconv.Atoi(replyNumString)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error converting reply number to int: %v", err), http.StatusBadRequest)
		return
	}
	var data ReplyRequest
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		s.returnError(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}
	s.redisClient.SaveReply(userId1, userId2, questionNum, replyNum, data.Reply)
	w.Write([]byte("Reply saved!"))
}

func (s *Server) GetRandomQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	questions := questions.GetRandomQuestions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (s *Server) GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	questions := questions.GetQuestions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
