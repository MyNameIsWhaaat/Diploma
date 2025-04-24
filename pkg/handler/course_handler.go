package handler

import (
	"net/http"
	"strconv"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type getAllCoursesResponse struct {
	Data []models.CourseWithProgress `json:"data"`
}

func (h *Handler) getAllCourseWithProgress(c *gin.Context) {
	userId, err := getUserId(c)
	courses, err := h.services.Course.GetCourseWithProgress(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("failed to get courses: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCoursesResponse{
		Data: courses,
	})
}

func (h *Handler) getAllCourseWithoutProgress(c *gin.Context) {
	userId, err := getUserId(c)
	courses, err := h.services.Course.GetCourseWithoutProgress(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("failed to get courses: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCoursesResponse{
		Data: courses,
	})
}

func (h *Handler) startCourse(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid course id")
		return
	}

	err = h.services.Course.StartCourse(userId, courseId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}