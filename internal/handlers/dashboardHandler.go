package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetRecentlyAssignedTasks(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Println("Invalid limit parameter:", err)
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	tasks, err := h.dashboardService.GetRecentlyAssignedTasks(user.UID, limit)
	if err != nil {
		log.Println("Failed to get recently assigned tasks:", err)
		http.Error(w, "Failed to get recently assigned tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("Failed to encode tasks:", err)
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}

func (h *DashboardHandler) GetInProgressTasks(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Println("Invalid limit parameter:", err)
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	tasks, err := h.dashboardService.GetInProgressTasks(user.UID, limit)
	if err != nil {
		log.Println("Failed to get in-progress tasks:", err)
		http.Error(w, "Failed to get in-progress tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("Failed to encode tasks:", err)
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}

func (h *DashboardHandler) GetApproachingDeadlineTasks(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Println("Invalid limit parameter:", err)
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	tasks, err := h.dashboardService.GetApproachingDeadlineTasks(user.UID, limit)
	if err != nil {
		log.Println("Failed to get approaching deadline tasks:", err)
		http.Error(w, "Failed to get approaching deadline tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("Failed to encode tasks:", err)
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}
