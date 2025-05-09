package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDTO models.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		log.Println("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateUser(&userDTO); err != nil {
		log.Println("Failed to register user: ", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userDTO *models.UserResponseDTO
	userDTO, err = h.userService.GetUserByFirebaseUID(user.UID)

	if err != nil {
		log.Println("Failed to get user:", err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDTO)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userEmail, ok := user.Claims["email"].(string)
	if !ok || userEmail == "" {
		log.Println("Email not found in Firebase token claims")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sq struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sq); err != nil {
		log.Println("Invalid request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userResponsesDTO, err := h.userService.SearchUsers(sq.Query)
	if err != nil {
		log.Println("Error searching users:", err)
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	filteredUsers := []models.UserResponseDTO{}
	for _, u := range userResponsesDTO {
		if u.Email != userEmail { // Compare with extracted email
			filteredUsers = append(filteredUsers, u)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredUsers)
}
