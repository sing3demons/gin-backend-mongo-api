package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sing3demons/gin-backend-api/db"
	"github.com/sing3demons/gin-backend-api/models"
	"github.com/sing3demons/gin-backend-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	db *mongo.Database
}

func New(db *mongo.Database) *handler {
	return &handler{db}
}

func (h *handler) Collection() *mongo.Collection {
	return h.db.Collection("users")
}

func (h *handler) GetProfile(c *gin.Context) {
	sub, _ := c.Get("sub")

	id, err := primitive.ObjectIDFromHex(sub.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	db := db.New(h.Collection())
	user, err := db.GetById(bson.D{{Key: "_id", Value: id}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, user)
}

func (h *handler) GetUsers(c *gin.Context) {
	db := db.New(h.Collection())

	filter := bson.D{{}}

	users, err := db.GetList(filter)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, users)
}

type Register struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *handler) Register(c *gin.Context) {
	var body Register

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Fullname:  body.Fullname,
		Password:  string(bcryptPassword),
		Email:     body.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Blogs:     []models.Blog{},
	}

	db := db.New(h.Collection())
	uExit := db.CheckEmail(bson.D{{Key: "email", Value: user.Email}})

	if len(uExit.Email) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already exists",
		})
		return
	}

	result, err := db.Create(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusCreated, gin.H{"user_id": result.InsertedID})
}

func (h handler) Login(c *gin.Context) {
	var body Login
	c.ShouldBindJSON(&body)

	db := db.New(h.Collection())

	exit, err := db.Search(bson.D{{Key: "email", Value: body.Email}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(exit.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   exit.ID.Hex(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, gin.H{"access_token": tokenString})
}

func (h *handler) GetUser(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	db := db.New(h.Collection())

	user, err := db.GetById(bson.D{{Key: "_id", Value: userId}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	utils.ResponseJsonWithLogger(c, http.StatusOK, user)
}
