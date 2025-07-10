package routes

import (
	"taskive/controllers"
	"taskive/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controllers.AuthController,
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
	commentController *controllers.CommentController,
) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(middlewares.LoggerMiddleware())

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
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