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
	Status string
}

type TaskRepository interface {
	CreateTask(task *Task) error
	GetTaskByID(id uint) (*Task, error)
	GetTasksByUserID(userID uint, filter TaskFilter) ([]Task, error)
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

func (r *taskRepository) GetTasksByUserID(userID uint, filter TaskFilter) ([]Task, error) {
	var tasks []Task
	query := r.db.Preload("User").Where("user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
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