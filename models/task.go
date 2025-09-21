package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'pending'" json:"status"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
)

type TaskFilter struct {
	Status  string
	Page    int
	Limit   int
	SortBy  string // created_at, updated_at, title, status
	Order   string // asc, desc
}

type TasksResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

type TaskRepository interface {
	CreateTask(task *Task) error
	GetTaskByID(id uint) (*Task, error)
	GetTasksByUserID(userID uint, filter TaskFilter) (*TasksResponse, error)
	UpdateTask(task *Task) error
	DeleteTask(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return err
	}
	// Reload the task with user data
	return r.db.Preload("User").First(task, task.ID).Error
}

func (r *taskRepository) GetTaskByID(id uint) (*Task, error) {
	var task Task
	err := r.db.Preload("User").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetTasksByUserID(userID uint, filter TaskFilter) (*TasksResponse, error) {
	var tasks []Task
	var total int64

	// Base query
	baseQuery := r.db.Model(&Task{}).Where("user_id = ?", userID)

	// Apply status filter
	if filter.Status != "" {
		baseQuery = baseQuery.Where("status = ?", filter.Status)
	}

	// Count total records
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Set defaults for pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Max limit
	}

	// Calculate offset
	offset := (filter.Page - 1) * filter.Limit

	// Build query with preload
	query := r.db.Preload("User").Where("user_id = ?", userID)

	// Apply status filter
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Apply sorting
	orderBy := "created_at DESC" // default
	if filter.SortBy != "" {
		validSortFields := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"title":      true,
			"status":     true,
		}
		if validSortFields[filter.SortBy] {
			order := "DESC"
			if filter.Order == "asc" {
				order = "ASC"
			}
			orderBy = filter.SortBy + " " + order
		}
	}

	// Apply pagination and sorting
	err := query.Order(orderBy).Limit(filter.Limit).Offset(offset).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &TasksResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *taskRepository) UpdateTask(task *Task) error {
	if err := r.db.Save(task).Error; err != nil {
		return err
	}
	// Reload the task with user data
	return r.db.Preload("User").First(task, task.ID).Error
}

func (r *taskRepository) DeleteTask(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}