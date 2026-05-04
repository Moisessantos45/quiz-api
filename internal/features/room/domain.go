package room

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"quiz/internal/shared/models"
	"strings"
)

type CreateRoomResult struct {
	Game      *models.Game      `json:"game"`
	User      *models.User      `json:"user"`
	Departure *models.Departure `json:"departure"`
}

type JoinRoomResult struct {
	Game      *models.Game      `json:"game"`
	User      *models.User      `json:"user"`
	Departure *models.Departure `json:"departure"`
}

type RoomRepository interface {
	CreateGame(ctx context.Context, game *models.Game) error
	GetGameByKey(ctx context.Context, key string) (*models.Game, error)
	GetGameByID(ctx context.Context, id uint64) (*models.Game, error)
	UpdateGameStatus(ctx context.Context, gameID uint64, status string) error
	CreateUser(ctx context.Context, user *models.User) error
	CreateDeparture(ctx context.Context, departure *models.Departure) error
}

type RoomService interface {
	CreateRoom(ctx context.Context, name string, nickname string, avatarURL string, categoryID uint64) (*CreateRoomResult, error)
	JoinRoom(ctx context.Context, key string, nickname string, avatarURL string) (*JoinRoomResult, error)
	StartGame(ctx context.Context, gameID uint64) error
}

func generateKey() (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, 6)
	for i := range key {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		key[i] = chars[n.Int64()]
	}
	return string(key), nil
}

func newUser(nickname string, avatarURL string) (*models.User, error) {
	name := strings.TrimSpace(nickname)
	if name == "" {
		return nil, fmt.Errorf("el nickname no puede estar vacio")
	}
	return &models.User{Nickname: name, AvatarURL: avatarURL}, nil
}
