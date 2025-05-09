package handler

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/service"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler{
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.RedirectTrailingSlash = false 

	// CORS Middleware с разрешением всех источников
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
			courses.POST("/start/:id", h.startCourse)

			levels := courses.Group("/:course_id/levels")
			{
				levels.GET("", h.getCourseLevels)
				
			}
			courses.GET("/level/:id", h.getTasksByLevel)
			courses.POST("/complete/level/:id", h.completeLevel)
		}
		api.POST("/submitAnswer", h.submitAnswerHandler)
	}

	return router
}
