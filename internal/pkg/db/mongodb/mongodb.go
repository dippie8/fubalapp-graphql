package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

var Db *mongo.Client

type DbConfig struct {
	Uri string `yaml:"Uri"`
}

func InitDB() {
	dbConfig := DbConfig{}
	yamlFile, err := ioutil.ReadFile("parameters.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &dbConfig)

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConfig.Uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	Db = client
	//defer client.Disconnect(ctx)
}