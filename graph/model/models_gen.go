// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type DeleteGame struct {
	ID string `json:"id"`
}

type Game struct {
	ID        	string `bson:"_id"`
	Player1     *Player `json:"player1"`
	Player2     *Player `json:"player2"`
	Player3     *Player `json:"player3"`
	Player4     *Player `json:"player4"`
	Score12     int    `json:"score12"`
	Score34     int    `json:"score34"`
	CreatedBy   string `json:"createdBy"`
	DeltaPoints int    `json:"deltaPoints"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewGame struct {
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
	Player4 string `json:"player4"`
	Score12 int    `json:"score12"`
	Score34 int    `json:"score34"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Player struct {
	Username     string      `bson:"_id"`
	CareerWin    int         `json:"careerWin"`
	CareerPlayed int         `json:"careerPlayed"`
	GoldMedals   int         `json:"goldMedals"`
	SilverMedals int         `json:"silverMedals"`
	BronzeMedals int         `json:"bronzeMedals"`
	Color        string      `json:"color"`
	IsAdmin      int         `json:"isAdmin"`
	Teammates    []*Teammate `json:"teammates"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type Standing struct {
	Username string `bson:"_id"`
	Win      int    `json:"win"`
	Played   int    `json:"played"`
	Elo      int    `json:"elo"`
}

type Teammate struct {
	Username      string `json:"username"`
	GamesTogether int    `json:"gamesTogether"`
	WinTogether   int    `json:"winTogether"`
	GamesAgainst  int    `json:"gamesAgainst"`
	WinAgainst    int    `json:"winAgainst"`
}
