package handler

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/service"
	"github.com/gin-gonic/gin"
)

type ServiceInterface interface {
	CreateUser(user domain.User) (int, error)
}

type Handler struct {
	services *service.Service

	serviceV2 ServiceInterface
}

func NewHandler(services *service.Service, serviceV2 ServiceInterface) *Handler {
	return &Handler{services: services, serviceV2: serviceV2}
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
		courses := api.Group("/courses")
		{
			courses.GET("/with-progress", h.getAllCourseWithProgress)
			courses.GET("/without-progress", h.getAllCourseWithoutProgress)
			courses.POST("/:id/start", h.startCourse)

			levels := courses.Group("/:id/levels")
			{
				levels.GET("", h.getCourseLevels)
				levels.POST("/:id/complete", h.completeLevel)
			}
		}

	}

	return router
}
