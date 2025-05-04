package handler

import (
	"net/http"
	"strconv"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type getCoursesWithProgressResponse struct {
	Data []domain.CourseWithProgress `json:"data"`
}

type getCoursesWithoutProgressResponse struct {
	Data []domain.CourseWithoutProgress `json:"data"`
}

func (h *Handler) getAllCourseWithProgress(c *gin.Context) {
	userId, err := getUserId(c)
	courses, err := h.services.GetCourseWithProgress(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("failed to get courses: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, getCoursesWithProgressResponse{
		Data: courses,
	})
}

func (h *Handler) getAllCourseWithoutProgress(c *gin.Context) {
	userId, err := getUserId(c)
	courses, err := h.services.GetCourseWithoutProgress(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("failed to get courses: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, getCoursesWithoutProgressResponse{
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

	err = h.services.StartCourse(userId, courseId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}