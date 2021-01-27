package players

import (
	"context"
	database "fubalapp-graphql/internal/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Teammate data
type Teammate struct {
	Username      string `json:"username"`
	GamesTogether int    `json:"gamesTogether"`
	WinTogether   int    `json:"winTogether"`
	GamesAgainst  int    `json:"gamesAgainst"`
	WinAgainst    int    `json:"winAgainst"`
}

//Player data
type Player struct {
	Username     string      `json:"username" bson:"_id"`
	CareerWin    int         `json:"careerWin"`
	CareerPlayed int         `json:"careerPlayed"`
	GoldMedals   int         `json:"goldMedals"`
	SilverMedals int         `json:"silverMedals"`
	BronzeMedals int         `json:"bronzeMedals"`
	Color        string      `json:"color"`
	IsAdmin      int         `json:"isAdmin"`
	Teammates    []*Teammate `json:"teammates"`
}

// CreateFromUsername creates a player starting from a username
func CreateFromUsername(username string) (string, error) {
	var player = &Player{
		Username:     username,
		CareerWin:    0,
		CareerPlayed: 0,
		GoldMedals:   0,
		SilverMedals: 0,
		Color:        "none",
		IsAdmin:      0,
		Teammates:    []*Teammate{},
	}
	collection := database.Db.Database("qlsr").Collection("players")
	_, err := collection.InsertOne(context.TODO(), player)

	return player.Username, err
}

// Update data about users when new games are played
func Update(winners [2]string, losers [2]string) error {
	collection := database.Db.Database("qlsr").Collection("players")
	addWin := func(usr string) error {
		_, err := collection.UpdateOne(
			context.TODO(), bson.M{"_id": usr},
			bson.D{
				{"$inc", bson.D{{"careerplayed", 1}}},
				{"$inc", bson.D{{"careerwin", 1}}},
			},
			options.Update().SetUpsert(true),
		)
		return err
	}
	addDefeat := func(usr string) error {
		_, err := collection.UpdateOne(
			context.TODO(), bson.M{"_id": usr},
			bson.D{
				{"$inc", bson.D{{"careerplayed", 1}}},
			},
			options.Update().SetUpsert(true),
		)
		return err
	}
	err := addWin(winners[0])
	err = addWin(winners[1])
	err = addDefeat(losers[0])
	err = addDefeat(losers[1])

	winner1, err := Get(winners[0])
	winner2, err := Get(winners[1])
	loser1, err := Get(losers[0])
	loser2, err := Get(losers[1])
	if err != nil {
		log.Panic(err)
	}

	err = winner1.updateWinner(winner2, loser1, loser2)
	err = winner2.updateWinner(winner1, loser1, loser2)
	err = loser1.updateLoser(loser2, winner1, winner2)
	err = loser2.updateLoser(loser1, winner1, winner2)

	//
	return err
}

func (player Player) updateWinner(teammate, opponent1, opponent2 *Player) error {
	var winnerTeammates []*Teammate
	foundTeammate := false
	foundOpponent1 := false
	foundOpponent2 := false
	collection := database.Db.Database("qlsr").Collection("players")

	for _, tm := range player.Teammates {
		if tm.Username == teammate.Username {
			foundTeammate = true
			winnerTeammates = append(
				winnerTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether + 1,
					WinTogether:   tm.WinTogether + 1,
					GamesAgainst:  tm.GamesAgainst,
					WinAgainst:    tm.WinAgainst,
				},
			)
		} else if tm.Username == opponent1.Username {
			foundOpponent1 = true
			winnerTeammates = append(
				winnerTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst + 1,
					WinAgainst:    tm.WinAgainst + 1,
				},
			)
		} else if tm.Username == opponent2.Username {
			foundOpponent2 = true
			winnerTeammates = append(
				winnerTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst + 1,
					WinAgainst:    tm.WinAgainst + 1,
				},
			)
		} else {
			winnerTeammates = append(
				winnerTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst,
					WinAgainst:    tm.WinAgainst,
				},
			)
		}
	}
	if !foundTeammate {
		winnerTeammates = append(
			winnerTeammates,
			&Teammate{
				Username:      teammate.Username,
				GamesTogether: 1,
				WinTogether:   1,
				GamesAgainst:  0,
				WinAgainst:    0,
			},
		)
	}
	if !foundOpponent1 {
		winnerTeammates = append(
			winnerTeammates,
			&Teammate{
				Username:      opponent1.Username,
				GamesTogether: 0,
				WinTogether:   0,
				GamesAgainst:  1,
				WinAgainst:    1,
			},
		)
	}
	if !foundOpponent2 {
		winnerTeammates = append(
			winnerTeammates,
			&Teammate{
				Username:      opponent2.Username,
				GamesTogether: 0,
				WinTogether:   0,
				GamesAgainst:  1,
				WinAgainst:    1,
			},
		)
	}
	_, err := collection.UpdateOne(
		context.TODO(), bson.M{"_id": player.Username},
		bson.D{
			{"$set", bson.D{{"teammates", winnerTeammates}}},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (player Player) updateLoser(teammate, opponent1, opponent2 *Player) error {
	var loserTeammates []*Teammate
	foundTeammate := false
	foundOpponent1 := false
	foundOpponent2 := false
	collection := database.Db.Database("qlsr").Collection("players")

	for _, tm := range player.Teammates {
		if tm.Username == teammate.Username {
			foundTeammate = true
			loserTeammates = append(
				loserTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether + 1,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst,
					WinAgainst:    tm.WinAgainst,
				},
			)
		} else if tm.Username == opponent1.Username {
			foundOpponent1 = true
			loserTeammates = append(
				loserTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst + 1,
					WinAgainst:    tm.WinAgainst,
				},
			)
		} else if tm.Username == opponent2.Username {
			foundOpponent2 = true
			loserTeammates = append(
				loserTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst + 1,
					WinAgainst:    tm.WinAgainst,
				},
			)
		} else {
			loserTeammates = append(
				loserTeammates,
				&Teammate{
					Username:      tm.Username,
					GamesTogether: tm.GamesTogether,
					WinTogether:   tm.WinTogether,
					GamesAgainst:  tm.GamesAgainst,
					WinAgainst:    tm.WinAgainst,
				},
			)
		}
	}
	if !foundTeammate {
		loserTeammates = append(
			loserTeammates,
			&Teammate{
				Username:      teammate.Username,
				GamesTogether: 1,
				WinTogether:   0,
				GamesAgainst:  0,
				WinAgainst:    0,
			},
		)
	}
	if !foundOpponent1 {
		loserTeammates = append(
			loserTeammates,
			&Teammate{
				Username:      opponent1.Username,
				GamesTogether: 0,
				WinTogether:   0,
				GamesAgainst:  1,
				WinAgainst:    0,
			},
		)
	}
	if !foundOpponent2 {
		loserTeammates = append(
			loserTeammates,
			&Teammate{
				Username:      opponent2.Username,
				GamesTogether: 0,
				WinTogether:   0,
				GamesAgainst:  1,
				WinAgainst:    0,
			},
		)
	}
	_, err := collection.UpdateOne(
		context.TODO(), bson.M{"_id": player.Username},
		bson.D{
			{"$set", bson.D{{"teammates", loserTeammates}}},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// Get data about a specific user
func Get(username string) (*Player, error) {
	collection := database.Db.Database("qlsr").Collection("players")
	filter := bson.M{"_id": username}
	var player *Player
	err := collection.FindOne(context.TODO(), filter).Decode(&player)
	return player, err
}

// GetAll return data about every player
func GetAll() []*Player {
	collection := database.Db.Database("qlsr").Collection("players")

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"goldmedals", -1}, {"silvermedals", -1}, {"bronzemedals", -1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	var players []*Player

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var p *Player

		err = cursor.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}

		players = append(players, p)
	}
	return players
}
