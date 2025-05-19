package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubmitAnswerRequest struct {
	TaskID int         `json:"task_id"`
	Answer interface{} `json:"answer"`
}

type SubmitAnswerResponse struct {
	IsCorrect bool `json:"is_correct"`
	XPEarned  int  `json:"xp_earned"`
}

func (h *Handler) submitAnswerHandler(c *gin.Context) {
	var req SubmitAnswerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	isCorrect, xp, err := h.service.SubmitAnswer(userId, req.TaskID, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("submit failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, SubmitAnswerResponse{
		IsCorrect: isCorrect,
		XPEarned:  xp,
	})
}

func (h *Handler) getTasksByLevel(c *gin.Context) {
	
	levelIdParam := c.Param("id")
	levelId, err := strconv.Atoi(levelIdParam)
	if err != nil || levelId <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid level id")
		return
	}

	tasks, err := h.service.GetTasksByLevel(levelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) getTaskVariants(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	variants, err := h.service.GetTaskVariants(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task variants"})
		return
	}

	c.JSON(http.StatusOK, variants)
}

func (h *Handler) getReviewTasks(c *gin.Context) {
	
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	levelIdParam := c.Param("level_id")
	levelId, err := strconv.Atoi(levelIdParam)
	if err != nil || levelId <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid level id")
		return
	}
	tasks, err := h.service.GetReviewTasks(userId, levelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) getMatchPairs(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	pairs, err := h.service.GetMatchPairsByTaskID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get match pairs"})
		return
	}

	c.JSON(http.StatusOK, pairs)
}