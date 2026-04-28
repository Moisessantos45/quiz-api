package user

import (
	"errors"
	"quiz/internal/shared/models"
)

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) UserService {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Register(nickname string, gameID uint64) (*UserRegistrationResult, error) {
	user, err := NewUser(nickname)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	departure, err := NewDeparture(user.ID, gameID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.CreateDeparture(departure); err != nil {
		return nil, err
	}

	return &UserRegistrationResult{User: user, Departure: departure}, nil
}

func (uc *UserUseCase) GetByID(id uint64) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("El ID del usuario no puede ser cero")
	}

	return uc.repo.GetByID(id)
}
