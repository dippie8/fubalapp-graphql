package games

import (
	"context"
	"fmt"
	database "fubalapp-graphql/internal/pkg/db/mongodb"
	"fubalapp-graphql/internal/players"
	"fubalapp-graphql/internal/standings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// #1
type Game struct {
	ID      	string 	`bson:"_id"`
	Player1 	string	`json:"player1"`
	Player2 	string	`json:"player2"`
	Player3 	string	`json:"player3"`
	Player4 	string	`json:"player4"`
	Score12 	int    	`json:"score12"`
	Score34 	int    	`json:"score34"`
	CreatedBy	string	`json:"createdBy"`
}

func (game Game) Save() string {
	collection := database.Db.Database("qlsr").Collection("games")
	insertResult, err := collection.InsertOne(context.TODO(), game)

	if err != nil {
		log.Fatal(err)
	}

	// update classifica
	var winners [2]string
	var losers [2]string

	if game.Score12 > game.Score34 {
		winners = [2]string{game.Player1, game.Player2}
		losers = [2]string{game.Player3, game.Player4}
	} else if game.Score34 > game.Score12 {
		winners = [2]string{game.Player3, game.Player4}
		losers = [2]string{game.Player1, game.Player2}
	}

	err = standings.Update(winners, losers)

	if err != nil {
		log.Fatal(err)
	}

	// update statistiche giocatori
	err = players.Update(winners, losers)
	if err != nil {
		log.Fatal(err)
	}


	id := fmt.Sprintf("%s", insertResult.InsertedID)
	return id
}

func GetAll() []*Game {
	collection := database.Db.Database("qlsr").Collection("games")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	var games []*Game

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var g *Game

		err = cursor.Decode(&g)
		if err != nil {
			log.Fatal(err)
		}

		games = append(games, g)
	}
	return games
}

func GetLatest(n int64) []*Game {
	collection := database.Db.Database("qlsr").Collection("games")
	filter := bson.D{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}}).SetLimit(n)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	var games []*Game

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var g *Game

		err = cursor.Decode(&g)
		if err != nil {
			log.Fatal(err)
		}

		games = append(games, g)
	}
	return games
}