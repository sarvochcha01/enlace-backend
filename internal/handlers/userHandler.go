package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
