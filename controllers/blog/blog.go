package blog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/db"
	"github.com/sing3demons/gin-backend-api/models"
	"github.com/sing3demons/gin-backend-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	db *mongo.Database
}

func New(db *mongo.Database) *handler {
	return &handler{db}
}

func (h *handler) Collection() *mongo.Collection {
	return h.db.Collection("blogs")
}

func (h *handler) GetBlogs(c *gin.Context) {
	db := db.NewBlog(h.Collection())

	filter := bson.D{{}}

	blogs, err := db.FindAll(filter)
	fmt.Println(blogs)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, gin.H{"blogs": blogs})
}

func (h *handler) CreateBlog(c *gin.Context) {
	type Request struct {
		Topic string `json:"topic"`
	}

	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	db := db.NewBlog(h.Collection())

	sub, _ := c.Get("sub")

	id, err := primitive.ObjectIDFromHex(sub.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	body := models.Blog{
		ID:        primitive.NewObjectID(),
		Topic:     req.Topic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    id,
	}
	result, err := db.CreateBlog(body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusCreated, gin.H{"result": result})
}

func (h *handler) GetBlog(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	db := db.NewBlog(h.Collection())

	blog, err := db.FindById(bson.D{{Key: "_id", Value: id}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, gin.H{"blog": blog})
}
