package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hrusfandi/sb-task-management/handlers"
	authMiddleware "github.com/hrusfandi/sb-task-management/middleware"
	"github.com/hrusfandi/sb-task-management/models"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userRepo := models.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)

	taskRepo := models.NewTaskRepository(db)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuth)

			r.Route("/tasks", func(r chi.Router) {
				r.Post("/", taskHandler.CreateTask)
				r.Get("/", taskHandler.ListTasks)
				r.Get("/{id}", taskHandler.GetTask)
				r.Put("/{id}", taskHandler.UpdateTask)
				r.Delete("/{id}", taskHandler.DeleteTask)
			})
		})
	})

	return r
}