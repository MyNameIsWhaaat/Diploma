package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)
type LevelWithUserProgressResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	IsCompleted bool       `json:"is_completed"`
	IsCurrent   bool       `json:"is_current"`
	XPEarned    int        `json:"xp_earned"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (h *Handler) getCourseLevels(c *gin.Context) {

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseId, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid course id")
		return
	}

	levels, err := h.service.GetCourseLevelsForUser(userId, courseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Преобразуем в безопасный ответ
	var response []LevelWithUserProgressResponse
	for _, l := range levels {
		response = append(response, LevelWithUserProgressResponse{
			ID:          l.ID,
			Title:       l.Title,
			IsCompleted: l.IsCompleted.Valid && l.IsCompleted.Bool,
			IsCurrent:   l.IsCurrent.Valid && l.IsCurrent.Bool,
			XPEarned:    int(l.XPEarned.Int64),
			StartedAt:   nullableTime(l.StartedAt),
			CompletedAt: nullableTime(l.CompletedAt),
		})
	}

	c.JSON(http.StatusOK, response)
}

func nullableTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func (h *Handler) completeLevel(c *gin.Context) {
	// Получение ID пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получение и валидация ID уровня из параметра URL
	levelIdParam := c.Param("id")
	levelId, err := strconv.Atoi(levelIdParam)
	if err != nil || levelId <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid level id")
		return
	}

	// Попытка завершить уровень
	err = h.service.CompleteLevel(userId, levelId)
	if err != nil {
		// Обработка ошибок бизнес-логики (например, уровень не найден, уже завершён, не принадлежит пользователю и т.п.)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}

func (h *Handler) startLevel(c *gin.Context) {
	// Получение ID пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получение и валидация ID уровня из параметра URL
	levelIdParam := c.Param("id")
	levelId, err := strconv.Atoi(levelIdParam)
	if err != nil || levelId <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid level id")
		return
	}

	// Попытка начать уровень
	err = h.service.StartLevel(userId, levelId)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not enrolled"):
			newErrorResponse(c, http.StatusForbidden, err.Error())
		case strings.Contains(err.Error(), "already started"):
			newErrorResponse(c, http.StatusConflict, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Успешный ответ
	c.JSON(http.StatusCreated, statusResponse{Status: "started"})
}