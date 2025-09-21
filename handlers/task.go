package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hrusfandi/sb-task-management/middleware"
	"github.com/hrusfandi/sb-task-management/models"
	"github.com/hrusfandi/sb-task-management/utils"
	"gorm.io/gorm"
)

type TaskHandler struct {
	taskRepo models.TaskRepository
}

func NewTaskHandler(taskRepo models.TaskRepository) *TaskHandler {
	return &TaskHandler{
		taskRepo: taskRepo,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		utils.RespondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Status == "" {
		req.Status = models.TaskStatusPending
	}

	if req.Status != models.TaskStatusPending &&
	   req.Status != models.TaskStatusInProgress &&
	   req.Status != models.TaskStatusCompleted {
		utils.RespondError(w, http.StatusBadRequest, "Invalid status value")
		return
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UserID:      userClaims.UserID,
	}

	if err := h.taskRepo.CreateTask(task); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	utils.RespondCreated(w, "Task created successfully", task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	page := 1
	limit := 10
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	filter := models.TaskFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  order,
	}

	// Validate status if provided
	if status != "" &&
	   status != models.TaskStatusPending &&
	   status != models.TaskStatusInProgress &&
	   status != models.TaskStatusCompleted {
		utils.RespondError(w, http.StatusBadRequest, "Invalid status value")
		return
	}

	result, err := h.taskRepo.GetTasksByUserID(userClaims.UserID, filter)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	utils.RespondSuccess(w, "Tasks fetched successfully", result)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskRepo.GetTaskByID(uint(taskID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Task not found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	if task.UserID != userClaims.UserID {
		utils.RespondError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondSuccess(w, "Task fetched successfully", task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskRepo.GetTaskByID(uint(taskID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Task not found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	if task.UserID != userClaims.UserID {
		utils.RespondError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title != "" {
		task.Title = strings.TrimSpace(req.Title)
		if task.Title == "" {
			utils.RespondError(w, http.StatusBadRequest, "Title cannot be empty")
			return
		}
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		if req.Status != models.TaskStatusPending &&
		   req.Status != models.TaskStatusInProgress &&
		   req.Status != models.TaskStatusCompleted {
			utils.RespondError(w, http.StatusBadRequest, "Invalid status value")
			return
		}
		task.Status = req.Status
	}

	if err := h.taskRepo.UpdateTask(task); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	utils.RespondSuccess(w, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskRepo.GetTaskByID(uint(taskID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusNotFound, "Task not found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	if task.UserID != userClaims.UserID {
		utils.RespondError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.taskRepo.DeleteTask(uint(taskID)); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	utils.RespondSuccess(w, "Task deleted successfully", nil)
}