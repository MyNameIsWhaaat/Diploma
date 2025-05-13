package handler

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	//"github.com/MyNameIsWhaaat/algo-learning/pkg/service"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

type Handler struct {
	//services *service.Service
	service ServiceInterface
}

type ServiceInterface interface {
	CreateUser(user domain.User) (int, error)
	SubmitAnswer(userID, taskID int, rawAnswer interface{}) (bool, int, error)
	GenerateToken(username, password string) (string, error)
	GetCourseWithoutProgress(userId int) ([]domain.CourseWithoutProgress, error)
	GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error)
	StartCourse(userID, courseID int) error
	CompleteLevel(userID, levelID int) (*domain.LevelResult, error)
	GetCourseLevelsForUser(userID, courseID int) ([]domain.LevelWithUserProgress, error)
	GetTasksByLevel(levelID int) ([]domain.Task, error)
	ParseToken(accessToken string) (int, error)
	StartLevel(userID, levelID int) error
	GetUserProfileData(userID int) (domain.UserProfileData, error)
	GetTaskVariants(taskID int) ([]domain.TaskVariant, error)
	GetReviewTasks(userID, levelID int) ([]domain.Task, error)
	
	GetLevelsWithAccess(userID, courseID int) ([]domain.LevelWithAccess, error)
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
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
		courses := api.Group("/courses")
		{
			courses.GET("/with-progress", h.getAllCourseWithProgress)
			courses.GET("/without-progress", h.getAllCourseWithoutProgress)
			courses.POST("/start/:id", h.startCourse)

			levels := courses.Group("/:course_id/levels")
			{
				levels.GET("", h.getCourseLevels)
				levels.GET("/levels-with-access", h.getLevelsWithAccess)

			}
			courses.GET("/level/:id/tasks", h.getTasksByLevel)
			courses.POST("/level/:id/start", h.startLevel)
			courses.POST("/complete/level/:id", h.completeLevel)
		}
		tasks := api.Group("/tasks")
		{
			tasks.GET("/:id/variants", h.getTaskVariants)
			tasks.GET("/review-tasks/:level_id", h.getReviewTasks)

		}
		api.POST("/submit_answer", h.submitAnswerHandler)
		api.GET("/user/me", h.getUserProfileData)
	}

	return router
}
