package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sukhmai/fun/pkg/db"
)

func (s *Server) UploadProfilePhotoHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := s.profileHandler(w, r)
	if err != nil {
		return
	}
	var data struct {
		Photo string `json:"photo"`
	}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	photoData, err := base64.StdEncoding.DecodeString(data.Photo)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error decoding base64 photo data: %v", err), http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("photos/%s-photo.jpg", userId)

	err = os.WriteFile(fileName, photoData, 0644)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error saving photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Photo uploaded successfully")
}

func (s *Server) GetProfilePhotoHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := s.profileHandler(w, r)
	if err != nil {
		return
	}
	fileName := fmt.Sprintf("photos/%s-photo.webp", userId)
	photoData, err := os.ReadFile(fileName)
	if err != nil {
		// TODO: pick a random placeholder photo and save that as the user's photo
		photoData, err = os.ReadFile("photos/placeholder.webp")
		if err != nil {
			s.returnError(w, fmt.Sprintf("error reading placeholder photo: %v", err), http.StatusInternalServerError)
			return
		}
	}

	encodedPhoto := base64.StdEncoding.EncodeToString(photoData)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"photo": encodedPhoto})
}

func (s *Server) GetDefaultPhotoHandler(w http.ResponseWriter, r *http.Request) {
	photoData, err := os.ReadFile("photos/placeholder.webp")
	if err != nil {
		s.returnError(w, fmt.Sprintf("error reading placeholder photo: %v", err), http.StatusInternalServerError)
		return
	}

	encodedPhoto := base64.StdEncoding.EncodeToString(photoData)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"photo": encodedPhoto})
}

func (s *Server) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	err = s.dbClient.SaveUser(r.Context(), user)
	if err != nil {
		s.returnError(w, fmt.Sprintf("error saving user to db: %v", err), http.StatusInternalServerError)
		return
	}

	response := Response{Message: "User saved!"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := s.returnError(w, "Authorization header is required", http.StatusUnauthorized)
		return "", err
	}

	// Bearer token is usually in the format "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		err := s.returnError(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return "", err
	}

	uuid := parts[1]
	return uuid, nil
}
