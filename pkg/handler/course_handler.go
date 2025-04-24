package handler

import (
	"net/http"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type getAllCoursesResponse struct {
	Data []models.CourseWithProgress `json:"data"`
}

func (h *Handler) getAllCourses(c *gin.Context) {
	courses, err := h.services.Course.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("failed to get courses: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCoursesResponse{
		Data: courses,
	})
}