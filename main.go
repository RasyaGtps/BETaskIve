package main

import (
	"fmt"
	"log"
	"taskive/config"
	"taskive/controllers"
	"taskive/middlewares"
	"taskive/models"
	"taskive/routes"
	"taskive/services"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize middleware with database connection
	middlewares.InitMiddleware(db)

	// Auto migrate database
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.Comment{},
		&models.Member{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db)
	projectService := services.NewProjectService(db)
	taskService := services.NewTaskService(db)
	commentService := services.NewCommentService(db)
	invitationService := services.NewInvitationService(db)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	projectController := controllers.NewProjectController(projectService)
	taskController := controllers.NewTaskController(taskService)
	commentController := controllers.NewCommentController(commentService)
	invitationController := controllers.NewInvitationController(invitationService)

	// Setup router
	router := routes.SetupRouter(
		authController,
		projectController,
		taskController,
		commentController,
		invitationController,
	)

	// Start server
	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 