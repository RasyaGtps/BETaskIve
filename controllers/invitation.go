package controllers

import (
	"net/http"
	"strconv"
	"taskive/services"

	"github.com/gin-gonic/gin"
)

type InvitationController struct {
	invitationService *services.InvitationService
}

func NewInvitationController(invitationService *services.InvitationService) *InvitationController {
	return &InvitationController{
		invitationService: invitationService,
	}
}

func (c *InvitationController) GetUserInvitations(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	invitations, err := c.invitationService.GetUserInvitations(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, invitations)
}

func (c *InvitationController) RespondToInvitation(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	var input struct {
		Accept bool `json:"accept" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.invitationService.RespondToInvitation(userID, uint(projectID), input.Accept); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
} 