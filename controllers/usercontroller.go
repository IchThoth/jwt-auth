package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ichthoth/jwt-auth/database"
	"github.com/ichthoth/jwt-auth/helpers"
	"github.com/ichthoth/jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate *validator.Validate = validator.New()

func PasswordHash(userpassword string) {

}

func VerifyPassword(userpassword string, inputPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(inputPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var users *models.User

		if err := c.BindJSON(&users); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(users)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": users.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking email"})
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": users.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking Phone number"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone numbeer already exists"})
		}

		users.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		users.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		users.ID = primitive.NewObjectID()
		users.User_id = users.ID.Hex()
		token, tokenRefresh, err := helpers.GenerateAllTokens(*users.Email, *users.First_name, *users.Last_name, *users.User_type, *&users.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		}
		users.Token = &token
		users.Refresh_token = &tokenRefresh

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, users)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var users models.User
		var foundUser models.User

		if err := c.BindJSON(&users); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": users.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is invalid"})
			return
		}
		passwordVerify, msg := VerifyPassword(*&users.Password, foundUser.Password)
		defer cancel()

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var users models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&users)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUsers() {

}
