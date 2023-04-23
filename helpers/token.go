package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ichthoth/jwt-auth/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var CurrentTime = time.Now().Local().Unix()

type SignDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_Type  string
	jwt.StandardClaims
}

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedTokenRefresh string, err error) {
	Claims := &SignDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_Type:  userType,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	RefreshClaims := &SignDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, Claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
	}
	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodES256, RefreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
	}

	return token, tokenRefresh, err
}

func UpdateAllTokens(signedtoken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedtoken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	uspsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &uspsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
	}

}

func ValidateTokens(Signedtoken string) (claims *SignDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		Signedtoken,
		&SignDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < CurrentTime {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
