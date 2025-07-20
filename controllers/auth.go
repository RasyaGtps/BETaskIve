package controllers

import (
	"net/http"
	"taskive/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AuthController struct {
	authService *services.AuthService
	validate    *validator.Validate
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var input services.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validate.Struct(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Register(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input services.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validate.Struct(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.authService.Login(input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	user, err := c.authService.GetCurrentUser(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *AuthController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	user, err := c.authService.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
} 