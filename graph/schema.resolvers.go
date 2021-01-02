package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"fubalapp-graphql/graph/generated"
	"fubalapp-graphql/graph/model"
	"fubalapp-graphql/internal/auth"
	"fubalapp-graphql/internal/games"
	"fubalapp-graphql/internal/players"
	"fubalapp-graphql/internal/standings"
	"fubalapp-graphql/internal/users"
	"fubalapp-graphql/pkg/jwt"
	"time"
)

func (r *mutationResolver) CreateGame(ctx context.Context, input model.NewGame) (*model.Game, error) {
	// authentication
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Game{}, fmt.Errorf("access denied")
	}

	var game games.Game
	dt := time.Now()

	game.ID = dt.String()
	game.Score12 = input.Score12
	game.Score34 = input.Score34
	game.Player1 = input.Player1
	game.Player2 = input.Player2
	game.Player3 = input.Player3
	game.Player4 = input.Player4
	game.CreatedBy = user.Username
	id, delta := game.Save()


	return &model.Game{
		ID:          id,
		Player1:     &model.Player{Username: game.Player1},
		Player2:     &model.Player{Username: game.Player2},
		Player3:     &model.Player{Username: game.Player3},
		Player4:     &model.Player{Username: game.Player4},
		Score12:     game.Score12,
		Score34:     game.Score34,
		CreatedBy:   game.CreatedBy,
		DeltaPoints: delta,
	}, nil
}

func (r *mutationResolver) DeleteGame(ctx context.Context, input model.DeleteGame) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	err := user.Create()
	if err != nil {
		return "", err
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Games(ctx context.Context, latest *int, player *string) ([]*model.Game, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.Game{}, fmt.Errorf("access denied")
	}

	var gamesResult []*model.Game
	var gameList []*games.Game

	nLatest := 0
	playerName := ""

	if latest != nil {
		nLatest = *latest
	}
	if player != nil {
		playerName = *player
	}

	gameList = games.Get(int64(nLatest), playerName)

	for _, game := range gameList {

		player1, _ := players.Get(game.Player1)
		player2, _ := players.Get(game.Player2)
		player3, _ := players.Get(game.Player3)
		player4, _ := players.Get(game.Player4)

		player1Model := &model.Player{
			Username:     player1.Username,
			CareerWin:    player1.CareerWin,
			CareerPlayed: player1.CareerPlayed,
			GoldMedals:   player1.GoldMedals,
			SilverMedals: player1.SilverMedals,
			BronzeMedals: player1.BronzeMedals,
			Color:        player1.Color,
			IsAdmin:      player1.IsAdmin,
			Teammates:    nil,
		}

		player2Model := &model.Player{
			Username:     player2.Username,
			CareerWin:    player2.CareerWin,
			CareerPlayed: player2.CareerPlayed,
			GoldMedals:   player2.GoldMedals,
			SilverMedals: player2.SilverMedals,
			BronzeMedals: player2.BronzeMedals,
			Color:        player2.Color,
			IsAdmin:      player2.IsAdmin,
			Teammates:    nil,
		}

		player3Model := &model.Player{
			Username:     player3.Username,
			CareerWin:    player3.CareerWin,
			CareerPlayed: player3.CareerPlayed,
			GoldMedals:   player3.GoldMedals,
			SilverMedals: player3.SilverMedals,
			BronzeMedals: player3.BronzeMedals,
			Color:        player3.Color,
			IsAdmin:      player3.IsAdmin,
			Teammates:    nil,
		}

		player4Model := &model.Player{
			Username:     player4.Username,
			CareerWin:    player4.CareerWin,
			CareerPlayed: player4.CareerPlayed,
			GoldMedals:   player4.GoldMedals,
			SilverMedals: player4.SilverMedals,
			BronzeMedals: player4.BronzeMedals,
			Color:        player4.Color,
			IsAdmin:      player4.IsAdmin,
			Teammates:    nil,
		}

		gamesResult = append(
			gamesResult,
			&model.Game{
				ID:          game.ID,
				Player1:     player1Model,
				Player2:     player2Model,
				Player3:     player3Model,
				Player4:     player4Model,
				Score12:     game.Score12,
				Score34:     game.Score34,
				CreatedBy:   game.CreatedBy,
				DeltaPoints: game.DeltaPoints,
			},
		)
	}
	return gamesResult, nil
}

func (r *queryResolver) Players(ctx context.Context, username *string) ([]*model.Player, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.Player{}, fmt.Errorf("access denied")
	}

	var playersList []*players.Player
	var playersResults []*model.Player

	if username != nil {
		pl, err := players.Get(*username)
		if err == nil {
			playersList = append(playersList, pl)
		}
	} else {
		playersList = players.GetAll()
	}

	for _, p := range playersList {
		var newPlayer *model.Player
		var teammates []*model.Teammate

		for _, tm := range p.Teammates {
			var teammate *model.Teammate
			teammate = &model.Teammate{
				Username:      tm.Username,
				GamesTogether: tm.GamesTogether,
				WinTogether:   tm.WinTogether,
				GamesAgainst:  tm.GamesAgainst,
				WinAgainst:    tm.WinAgainst,
			}
			teammates = append(teammates, teammate)
		}

		newPlayer = &model.Player{
			Username:     p.Username,
			CareerWin:    p.CareerWin,
			CareerPlayed: p.CareerPlayed,
			GoldMedals:   p.GoldMedals,
			SilverMedals: p.SilverMedals,
			BronzeMedals: p.BronzeMedals,
			Color:        p.Color,
			IsAdmin:      p.IsAdmin,
			Teammates:    teammates,
		}
		playersResults = append(playersResults, newPlayer)
	}
	return playersResults, nil
}

func (r *queryResolver) Standings(ctx context.Context) ([]*model.Standing, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return []*model.Standing{}, fmt.Errorf("access denied")
	}

	var standingResult []*model.Standing
	standingList := standings.GetAll()

	for _, standing := range standingList {
		standingResult = append(
			standingResult,
			&model.Standing{
				Username: standing.Username,
				Win:      standing.Win,
				Played:   standing.Played,
				Elo:      standing.Elo,
			},
		)
	}
	return standingResult, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Users(ctx context.Context) ([]*model.Player, error) {
	panic(fmt.Errorf("not implemented"))
}
