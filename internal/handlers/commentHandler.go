package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var CreateCommentDTO models.CreateCommentDTO
	var err error

	projectID := chi.URLParam(r, "projectID")

	CreateCommentDTO.ProjectID, err = uuid.Parse(projectID)
	if err != nil {
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	taskID := chi.URLParam(r, "taskID")

	CreateCommentDTO.TaskID, err = uuid.Parse(taskID)
	if err != nil {
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&CreateCommentDTO); err != nil {
		log.Fatal("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user *auth.Token
	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Fatal("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.commentService.CreateComment(&CreateCommentDTO, user.UID); err != nil {
		log.Fatal("Failed to create comment: ", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Comment created successfully"))

}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {

	commentID := chi.URLParam(r, "commentID")

	parsedCommentID, err := uuid.Parse(commentID)

	if err != nil {
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	_, err = middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Fatal("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	comment, err := h.commentService.GetComment(parsedCommentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)

}

func (h *CommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {

	var UpdateCommentDTO models.UpdateCommentDTO
	var err error

	projectID := chi.URLParam(r, "projectID")
	UpdateCommentDTO.ProjectID, err = uuid.Parse(projectID)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	commentID := chi.URLParam(r, "commentID")
	UpdateCommentDTO.CommentID, err = uuid.Parse(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&UpdateCommentDTO); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.commentService.UpdateComment(&UpdateCommentDTO, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment updated successfully"))
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var deleteCommentDTO models.DeleteCommentDTO
	var err error

	projectID := chi.URLParam(r, "projectID")
	deleteCommentDTO.ProjectID, err = uuid.Parse(projectID)
	if err != nil {
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	commentID := chi.URLParam(r, "commentID")
	deleteCommentDTO.CommentID, err = uuid.Parse(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.commentService.DeleteComment(&deleteCommentDTO, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment deleted successfully"))
}

func (h *CommentHandler) GetAllComments(w http.ResponseWriter, r *http.Request) {

}
