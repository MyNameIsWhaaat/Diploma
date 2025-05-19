package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubmitCodeRequest struct {
	TaskID   int    `json:"task_id"`
	UserCode string `json:"user_code"`
}

type SubmitCodeResponse struct {
	IsCorrect  bool `json:"is_correct"`
	Passed     int  `json:"passed"`
	Total      int  `json:"total"`
}

func (h *Handler) submitCodeHandler(c *gin.Context) {
	var req SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.CheckUserCode(c.Request.Context(), userID, req.TaskID, req.UserCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("check failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, SubmitCodeResponse{
		IsCorrect: result.IsCorrect,
		Passed:    result.Passed,
		Total:     result.Total,
	})
}
