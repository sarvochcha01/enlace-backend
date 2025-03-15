package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

type InvitationHandler struct {
	invitationService services.InvitationService
}

func NewInvitationHandler(is services.InvitationService) *InvitationHandler {
	return &InvitationHandler{invitationService: is}
}

func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var createInvitationDTO models.CreateInvitationDTO
	if err = json.NewDecoder(r.Body).Decode(&createInvitationDTO); err != nil {
		log.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err = h.invitationService.CreateInvitation(user.UID, &createInvitationDTO); err != nil {
		log.Println("Failed to create invitation:", err)
		http.Error(w, "Failed to create invitation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Invitation Created"))

}

func (h *InvitationHandler) GetInvitations(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invitations, err := h.invitationService.GetInvitations(user.UID)
	if err != nil {
		log.Println("Failed to get invitations:", err)
		http.Error(w, "Failed to get invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitations)
}

func (h *InvitationHandler) EditInvitation(w http.ResponseWriter, r *http.Request) {

	invitationID := chi.URLParam(r, "invitationID")
	parsedInvitationID, err := uuid.Parse(invitationID)
	if err != nil {
		log.Println("Invalid invitation id:", err)
		http.Error(w, "Invalid invitation id", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var EditInvitationRequest models.EditInvitationDTO

	if err = json.NewDecoder(r.Body).Decode(&EditInvitationRequest); err != nil {
		log.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	EditInvitationRequest.InvitationID = parsedInvitationID

	if err = h.invitationService.EditInvitation(user.UID, EditInvitationRequest); err != nil {
		log.Println("Failed to edit invitation:", err)
		http.Error(w, "Failed to edit invitation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Invitation edited successfully"))

}
