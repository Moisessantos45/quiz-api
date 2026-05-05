package user

import (
	"context"
	"fmt"
	"quiz/internal/shared/models"
	"strings"
)

type UserRegistrationResult struct {
	User      *models.User      `json:"user"`
	Departure *models.Departure `json:"departure"`
}

type UserRepository interface {
	GetByID(cxt context.Context, id uint64) (*models.User, error)
	Create(cxt context.Context, user *models.User) error
	CreateDeparture(cxt context.Context, departure *models.Departure) error
}

type UserService interface {
	Register(cxt context.Context, nickname string, avatarURL string, gameID uint64) (*UserRegistrationResult, error)
	GetByID(cxt context.Context, id uint64) (*models.User, error)
}

func NewUser(nickname string, avatarURL string) (*models.User, error) {
	name := strings.TrimSpace(nickname)
	if name == "" {
		return nil, fmt.Errorf("El nickname no puede estar vacio")
	}

	return &models.User{
		Nickname:  name,
		AvatarURL: avatarURL,
	}, nil
}

func NewDeparture(userID uint64, gameID uint64) (*models.Departure, error) {
	if userID == 0 {
		return nil, fmt.Errorf("El ID del usuario no puede ser cero")
	}
	if gameID == 0 {
		return nil, fmt.Errorf("El ID del juego no puede ser cero")
	}

	return &models.Departure{
		UserID: userID,
		GameID: gameID,
		Score:  0,
	}, nil
}
