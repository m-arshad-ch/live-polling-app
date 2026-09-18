package controllers

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"live-polling-app/models"
	"live-polling-app/utils"
)

type AuthController struct {
	db        *mongo.Database
	jwtSecret string
}

func NewAuthController(db *mongo.Database, jwtSecret string) *AuthController {
	return &AuthController{db: db, jwtSecret: jwtSecret}
}

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthController) Signup(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if len(req.Name) < 2 || len(req.Name) > 50 {
		c.JSON(400, gin.H{"error": "name must be between 2 and 50 characters"})
		return
	}
	if !strings.Contains(req.Email, "@") || len(req.Email) > 120 {
		c.JSON(400, gin.H{"error": "enter a valid email"})
		return
	}
	if len(req.Password) < 6 || len(req.Password) > 72 {
		c.JSON(400, gin.H{"error": "password must be 6 to 72 characters"})
		return
	}

	users := a.db.Collection("users")
	count, _ := users.CountDocuments(context.Background(), bson.M{"email": req.Email})
	if count > 0 {
		c.JSON(409, gin.H{"error": "email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create account"})
		return
	}

	user := models.User{
		Name: req.Name, Email: req.Email,
		PasswordHash: string(hash), CreatedAt: time.Now().Unix(),
	}

	result, err := users.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create account"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	token, err := utils.CreateToken(user.ID, a.jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(201, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID.Hex(), "name": user.Name, "email": user.Email},
	})
}

func (a *AuthController) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	err := a.db.Collection("users").FindOne(
		context.Background(), bson.M{"email": req.Email},
	).Decode(&user)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"error": "email or password is incorrect"})
		return
	}

	token, err := utils.CreateToken(user.ID, a.jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID.Hex(), "name": user.Name, "email": user.Email},
	})
}
