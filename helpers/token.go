package helpers

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ichthoth/jwt-auth/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

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
	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodES256, RefreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
	}

	return token, tokenRefresh, err
}
