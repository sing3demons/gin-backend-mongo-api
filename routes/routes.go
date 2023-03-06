package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/controllers/blog"
	"github.com/sing3demons/gin-backend-api/controllers/user"
	"github.com/sing3demons/gin-backend-api/logger"
	"github.com/sing3demons/gin-backend-api/middelwares"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewRouter(db *mongo.Database, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(logger.ZapLogger(log), logger.RecoveryWithZap(log, true))
	v1 := r.Group("/api/v1")

	controller := user.New(db)
	protect := middelwares.AuthJWT()

	auth := v1.Group("/users")
	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.GET("/", controller.GetUsers)
	auth.GET("/:id", controller.GetUser)
	auth.GET("/profile", protect, controller.GetProfile)

	blogs := v1.Group("/blogs")
	blogController := blog.New(db)

	blogs.POST("/", protect, blogController.CreateBlog)
	blogs.GET("/", blogController.GetBlogs)
	blogs.GET("/:id", blogController.GetBlog)

	return r
}
