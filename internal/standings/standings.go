package standings

import (
	"context"
	database "fubalapp-graphql/internal/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"sort"
)

const (
	eloConstant = 10
	minGamesEligibility = 5
)

// Standing data about a single player
type Standing struct {
	Username string `bson:"_id"`
	Win      int    `json:"win"`
	Played   int    `json:"played"`
	Elo      int    `json:"elo"`
}

type standingList []*Standing

func (s standingList) Len() int {
	return len(s)
}

func (s standingList) Less(i, j int) bool {

	var iPerc, jPerc float32

	if s[i].Played == 0 {
		iPerc = 0
	} else {
		iPerc = float32(s[i].Win) / float32(s[i].Played)
	}
	if s[i].Played >= minGamesEligibility {
		iPerc += 1
	}

	if s[j].Played == 0 {
		jPerc = 0
	} else {
		jPerc = float32(s[j].Win) / float32(s[j].Played)
	}
	if s[j].Played >= minGamesEligibility {
		jPerc += 1
	}


	return iPerc > jPerc
}

func (s standingList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// SubscribeUser creates a line in standings for a new user
func SubscribeUser(username string) (string, error) {
	var standing = &Standing{
		Username: username,
		Win:      0,
		Played:   0,
		Elo:      100,
	}
	collection := database.Db.Database("qlsr").Collection("standings")
	_, err := collection.InsertOne(context.TODO(), standing)

	return username, err
}

// Update standings when a new game is added
func Update(winners [2]string, losers [2]string) (int, error) {

	var winnersWinProbability float64

	w1, err := get(winners[0])
	w2, err := get(winners[1])
	l1, err := get(losers[0])
	l2, err := get(losers[1])
	if err != nil {
		log.Panic(err)
	}

	elo1 := math.Round(float64((w1.Elo + w2.Elo) / 2))
	elo2 := math.Round(float64((l1.Elo + l2.Elo) / 2))

	if elo1 < elo2 {
		winnersWinProbability = elo1 / elo2 * 0.5
	} else {
		winnersWinProbability = 1 - (elo2 / elo1 * 0.5)
	}

	delta := int(math.Round((1 - winnersWinProbability) * eloConstant))

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

func get(username string) (*Standing, error) {
	collection := database.Db.Database("qlsr").Collection("standings")
	filter := bson.M{"_id": username}
	var standing *Standing
	err := collection.FindOne(context.TODO(), filter).Decode(&standing)
	return standing, err
}

// GetAll returns the actual standings
func GetAll() []*Standing {
	collection := database.Db.Database("qlsr").Collection("standings")
	filter := bson.D{}
	// filter := bson.M{"played": {"$gte": 5}}
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

	sort.Sort(standingList(standings))

	return standings
}
