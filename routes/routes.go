package routes

import (
	"taskive/controllers"
	"taskive/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controllers.AuthController,
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
	commentController *controllers.CommentController,
	invitationController *controllers.InvitationController,
) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Middleware global
	router.Use(middlewares.LoggerMiddleware())

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/me", middlewares.AuthMiddleware(), authController.GetCurrentUser)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		// Users
		api.GET("/users", authController.GetUserByEmail)

		// Projects
		projects := api.Group("/projects")
		{
			projects.POST("", projectController.Create)
			projects.GET("", projectController.GetUserProjects)
			projects.GET("/:id", projectController.GetByID)
			projects.PUT("/:id", projectController.Update)
			projects.DELETE("/:id", projectController.Delete)
			projects.POST("/:id/invite", projectController.AddMember)

			// Tasks within project
			projects.GET("/:id/tasks", taskController.GetProjectTasks)
			projects.POST("/:id/tasks", taskController.Create)
		}

		// Invitations
		invitations := api.Group("/invitations")
		{
			invitations.GET("", projectController.GetUserInvitations)
			invitations.POST("/:id/accept", projectController.AcceptInvitation)
			invitations.POST("/:id/reject", projectController.RejectInvitation)
		}

		// Tasks
		tasks := api.Group("/tasks")
		{
			tasks.GET("/:id", taskController.GetByID)
			tasks.PUT("/:id", taskController.Update)
			tasks.DELETE("/:id", taskController.Delete)
			tasks.PATCH("/:id/status", taskController.UpdateStatus)

			// Comments within task
			tasks.GET("/:id/comments", commentController.GetTaskComments)
			tasks.POST("/:id/comments", commentController.Create)
		}

		// Comments
		comments := api.Group("/comments")
		{
			comments.DELETE("/:id", commentController.Delete)
		}
	}

	return router
} 