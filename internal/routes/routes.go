package routes

import (
	"database/sql"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/sarvochcha01/enlace-backend/internal/handlers"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/services"
	"github.com/sarvochcha01/enlace-backend/internal/websockets"
)

func SetupRoutes(r chi.Router, db *sql.DB, authClient *auth.Client) {

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	wsHub := websockets.NewWebSocketHub(authClient)
	wsHub.SetUserFinder(userService)
	go wsHub.Run()

	notificationRepository := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepository, wsHub, userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService, userService)

	projectMemberRepository := repositories.NewProjectMemberRepository(db)
	projectMemberService := services.NewProjectMemberService(projectMemberRepository, userService)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService)

	projectRepository := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepository, userService, projectMemberService)
	projectHandler := handlers.NewProjectHandler(projectService)

	taskRepository := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepository, userService, projectMemberService, notificationService)
	taskHandler := handlers.NewTaskHandler(taskService)

	commentRepository := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepository, userService, projectMemberService)
	commentHandler := handlers.NewCommentHandler(commentService)

	invitationRepository := repositories.NewInvitationRepository(db)
	invitationService := services.NewInvitationService(invitationRepository, userService, projectService, projectMemberService, notificationService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)

	dashboardRepository := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepository, userService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	authMiddleware := middlewares.NewAuthMiddleware(authClient)

	r.Route("/api/v1", func(api chi.Router) {

		api.Post("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Server Up"))
		})

		api.Route("/users", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)
			r.With(authMiddleware.FirebaseAuthMiddleware).Get("/", userHandler.GetUser)
			r.With(authMiddleware.FirebaseAuthMiddleware).Post("/search", userHandler.SearchUsers)

		})

		api.Route("/projects", func(r chi.Router) {
			r.Use(authMiddleware.FirebaseAuthMiddleware)
			r.Post("/", projectHandler.CreateProject)
			r.Get("/", projectHandler.GetAllProjectsForUser)

			r.Route("/{projectID}", func(r chi.Router) {
				r.Get("/", projectHandler.GetProjectByID)
				r.Put("/", projectHandler.UpdateProject)
				r.Delete("/", projectHandler.DeleteProject)

				r.Route("/project-members", func(r chi.Router) {
					// r.Get("/", projectMemberHandler.GetAllProjectMembers)
					r.Route("/{projectMemberID}", func(r chi.Router) {
						// r.Get("/", projectMemberHandler.GetProjectMember)
						r.Put("/", projectMemberHandler.UpdateProjectMember)
						// r.Delete("/", projectMemberHandler.DeleteProjectMember)
					})
				})

				r.Get("/join", projectHandler.GetProjectName)
				// TODO: Group join and leave, as well as updating the member roles (to be added) owner, editor, viewer into one handler func, such as projectHandler.UpdateMember or something
				r.Post("/join", projectHandler.JoinProject)
				r.Post("/leave", projectHandler.LeaveProject)

				r.Route("/tasks", func(r chi.Router) {
					r.Post("/", taskHandler.CreateTask)

					r.Route("/{taskID}", func(r chi.Router) {
						r.Get("/", taskHandler.GetTaskByID)
						r.Put("/", taskHandler.EditTask)
						r.Delete("/", taskHandler.DeleteTask)

						r.Route("/comments", func(r chi.Router) {
							r.Post("/", commentHandler.CreateComment)
							r.Get("/", commentHandler.GetAllCommentsForTask)
							r.Put("/{commentID}", commentHandler.EditComment)
							r.Get("/{commentID}", commentHandler.GetComment)
							r.Delete("/{commentID}", commentHandler.DeleteComment)
						})
					})

				})
			})
		})

		api.Route("/invitations", func(r chi.Router) {
			r.Use(authMiddleware.FirebaseAuthMiddleware)
			r.Post("/", invitationHandler.CreateInvitation)
			r.Get("/", invitationHandler.GetInvitations)

			r.Route("/{invitationID}", func(r chi.Router) {
				r.Put("/", invitationHandler.EditInvitation)
			})

			// To check if a user is invited to the project or not
			r.Get("/join-project/{projectID}", invitationHandler.HasInvitation)
		})

		api.Route("/dashboard", func(r chi.Router) {
			r.Use(authMiddleware.FirebaseAuthMiddleware)
			r.Get("/recently-assigned", dashboardHandler.GetRecentlyAssignedTasks)
			r.Get("/in-progress", dashboardHandler.GetInProgressTasks)
			r.Get("/approaching-deadline", dashboardHandler.GetApproachingDeadlineTasks)
			r.Get("/search", dashboardHandler.Search)
		})

		api.Route("/notifications", func(r chi.Router) {

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.FirebaseAuthMiddleware)
				r.Get("/", notificationHandler.GetAllNotificationsForUser)
				r.Post("/{notificationID}/read", notificationHandler.MarkNotificationAsRead)
			})

			r.Get("/ws", wsHub.HandleWebSocket)
		})

	})
}
