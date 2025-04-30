package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCourseLevels(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid course id")
		return
	}

	levels, err := h.services.GetByCourse(courseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, levels)
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
	err = h.services.CompleteLevel(userId, levelId)
	if err != nil {
		// Обработка ошибок бизнес-логики (например, уровень не найден, уже завершён, не принадлежит пользователю и т.п.)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}