package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserProfileData(c *gin.Context) {
	// Получение ID пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.service.GetUserProfileData(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}


	c.JSON(http.StatusOK, profile)
}

func (h *Handler) setUserNotNew(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.service.MarkUserNotNew(userID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}