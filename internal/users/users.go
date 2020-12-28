package users

import (
	"context"
	"fmt"
	database "fubalapp-graphql/internal/pkg/db/mongodb"
	"fubalapp-graphql/internal/players"
	"fubalapp-graphql/internal/standings"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"log"
)

type User struct {
	ID			string `bson:"_id"`
	Username	string `json:"name"`
	Password	string `json:"password"`
}

func (user *User) Create() error {

	hashedPassword, err := HashPassword(user.Password)

	collection := database.Db.Database("qlsr").Collection("users")
	_, err = collection.InsertOne(context.TODO(), User{ID: user.Username, Username: user.Username, Password: hashedPassword})
	if err != nil {
		return fmt.Errorf("user already exists")
	}
	_, err = players.CreateFromUsername(user.Username)
	if err != nil {
		return err
	}
	_, err = standings.SubscribeUser(user.Username)
	if err != nil {
		return err
	}

	return nil

}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (string, error) {
	collection := database.Db.Database("qlsr").Collection("users")
	filter := bson.M{"username": username}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if count != 1 {
		return "error", err
	}
	var user *User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Panic(err)
		return "error", err
	}
	return user.ID, nil
}

func (user *User) Authenticate() bool {
	var dbUser *User
	collection := database.Db.Database("qlsr").Collection("users")
	filter := bson.M{"username": user.Username}
	err := collection.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		return false
	}
	hashedPassword := dbUser.Password

	return CheckPasswordHash(user.Password, hashedPassword)
}
