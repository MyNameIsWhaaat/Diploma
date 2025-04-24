package handler

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler{
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		courses:= api.Group("/courses")
		{
			courses.GET("/with-progress", h.getAllCourseWithProgress)
			courses.GET("/without-progress", h.getAllCourseWithoutProgress)
		}
		
	}

	return router
}
