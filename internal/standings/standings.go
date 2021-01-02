package standings

import (
	"context"
	database "fubalapp-graphql/internal/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
)

type Standing struct {
	Username 	string	`bson:"_id"`
	Win 		int		`json:"win"`
	Played		int		`json:"played"`
	Elo			int		`json:"elo"`
}

func SubscribeUser(username string) (string, error) {
	var standing = &Standing{
		Username: username,
		Win: 0,
		Played: 0,
		Elo: 100,
	}
	collection := database.Db.Database("qlsr").Collection("standings")
	_, err := collection.InsertOne(context.TODO(), standing)

	return username, err
}

func Update (winners [2]string, losers [2]string) (int, error) {

	var winnersWinProbability float64
	const k = 10

	w1, err := Get(winners[0])
	w2, err := Get(winners[1])
	l1, err := Get(losers[0])
	l2, err := Get(losers[1])
	if err != nil {
		log.Panic(err)
	}

	elo1 := math.Round(float64((w1.Elo + w2.Elo) / 2))
	elo2 := math.Round(float64((l1.Elo + l2.Elo) / 2))

	if elo1 < elo2 {
		winnersWinProbability = elo1/elo2 * 0.5
	} else {
		winnersWinProbability = 1 - (elo2/elo1 * 0.5)
	}

	delta := int(math.Round((1 - winnersWinProbability) * k))

	collection := database.Db.Database("qlsr").Collection("standings")
	addWin := func(usr string) error {
		_, err := collection.UpdateOne(
			context.TODO(), bson.M{"_id": usr},
			bson.D{
				{"$inc", bson.D{{"win", 1}}},
				{"$inc", bson.D{{"played", 1}}},
				{"$inc", bson.D{{"elo", delta}}},
			},
			options.Update().SetUpsert(true),
		)
		return err
	}
	addLose := func(usr string) error {
		_, err := collection.UpdateOne(
			context.TODO(), bson.M{"_id": usr},
			bson.D{
				{"$inc", bson.D{{"win", 0}}},
				{"$inc", bson.D{{"played", 1}}},
				{"$inc", bson.D{{"elo", -delta}}},
			},
			options.Update().SetUpsert(true),
		)
		return err
	}

	err = addWin(winners[0])
	err = addWin(winners[1])
	err = addLose(losers[0])
	err = addLose(losers[1])

	return delta, err
}


func Get(username string) (*Standing, error) {
	collection := database.Db.Database("qlsr").Collection("standings")
	filter := bson.M{"_id": username}
	var standing *Standing
	err := collection.FindOne(context.TODO(), filter).Decode(&standing)
	return standing, err
}

func GetAll() []*Standing {
	collection := database.Db.Database("qlsr").Collection("standings")
	filter := bson.D{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"elo", -1}})
	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	var standings []*Standing

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var s *Standing

		err = cursor.Decode(&s)
		if err != nil {
			log.Fatal(err)
		}

		standings = append(standings, s)
	}
	return standings
}
