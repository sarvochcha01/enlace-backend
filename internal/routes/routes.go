package routes

import (
	"database/sql"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/sarvochcha01/enlace-backend/internal/handlers"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

func SetupRoutes(r chi.Router, db *sql.DB, authClient *auth.Client) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	projectMemberRepository := repositories.NewProjectMemberRepository(db)
	projectMemberService := services.NewProjectMemberService(projectMemberRepository, userService)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService)

	projectRepository := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepository, userService, projectMemberService)
	projectHandler := handlers.NewProjectHandler(projectService)

	taskRepository := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepository, userService, projectMemberService)
	taskHandler := handlers.NewTaskHandler(taskService)

	commentRepository := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepository, userService, projectMemberService)
	commentHandler := handlers.NewCommentHandler(commentService)

	authMiddleware := middlewares.NewAuthMiddleware(authClient)

	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/users", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)
		})

		api.Route("/projects", func(r chi.Router) {
			r.Use(authMiddleware.FirebaseAuthMiddleware)
			r.Post("/", projectHandler.CreateProject)
			r.Get("/", projectHandler.GetAllProjectsForUser)

			r.Post("/join/{projectID}", projectMemberHandler.CreateProjectMember)

			r.Route("/{projectID}", func(r chi.Router) {
				r.Get("/", projectHandler.GetProjectByID)

				r.Route("/tasks", func(r chi.Router) {
					r.Post("/", taskHandler.CreateTask)

					r.Route("/{taskID}", func(r chi.Router) {
						r.Route("/comments", func(r chi.Router) {
							r.Post("/", commentHandler.CreateComment)
							r.Get("/", commentHandler.GetAllComments)
							r.Put("/{commentID}", commentHandler.EditComment)
							r.Get("/{commentID}", commentHandler.GetComment)
							r.Delete("/{commentID}", commentHandler.DeleteComment)
						})
					})

				})
			})
		})

	})
}
