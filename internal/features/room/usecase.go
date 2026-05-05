package room

import (
	"context"
	"fmt"
	"quiz/internal/shared/models"
)

type RoomUseCase struct {
	repo RoomRepository
}

func NewRoomUseCase(repo RoomRepository) RoomService {
	return &RoomUseCase{repo: repo}
}

func (uc *RoomUseCase) CreateRoom(ctx context.Context, name string, nickname string, avatarURL string, categoryID uint64) (*CreateRoomResult, error) {
	user, err := newUser(nickname, avatarURL)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	key, err := generateKey()
	if err != nil {
		return nil, err
	}

	game := &models.Game{
		Key:        key,
		Name:       name,
		Status:     "waiting",
		UserID:     user.ID,
		CategoryID: categoryID,
	}
	if err := uc.repo.CreateGame(ctx, game); err != nil {
		return nil, err
	}

	departure := &models.Departure{
		UserID: user.ID,
		GameID: game.ID,
		Score:  0,
	}
	if err := uc.repo.CreateDeparture(ctx, departure); err != nil {
		return nil, err
	}

	return &CreateRoomResult{Game: game, User: user, Departure: departure}, nil
}

func (uc *RoomUseCase) JoinRoom(ctx context.Context, key string, nickname string, avatarURL string) (*JoinRoomResult, error) {
	game, err := uc.repo.GetGameByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if game.Status != "waiting" {
		return nil, fmt.Errorf("la partida ya inicio o finalizo")
	}

	user, err := newUser(nickname, avatarURL)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	departure := &models.Departure{
		UserID: user.ID,
		GameID: game.ID,
		Score:  0,
	}
	if err := uc.repo.CreateDeparture(ctx, departure); err != nil {
		return nil, err
	}

	players, _ := uc.repo.GetDeparturesByGameID(ctx, game.ID)

	return &JoinRoomResult{Game: game, User: user, Departure: departure, Players: players}, nil
}

func (uc *RoomUseCase) StartGame(ctx context.Context, gameID uint64) (*models.Game, error) {
	game, err := uc.repo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != "waiting" {
		return nil, fmt.Errorf("el juego no esta en espera")
	}
	if err := uc.repo.UpdateGameStatus(ctx, gameID, "started"); err != nil {
		return nil, err
	}
	return game, nil
}
