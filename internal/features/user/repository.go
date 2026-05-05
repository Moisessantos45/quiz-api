package user

import (
	"context"
	"quiz/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context,id uint64) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Create(ctx context.Context,user *models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error
}

func (r *PostgresRepository) CreateDeparture(ctx context.Context,departure *models.Departure) error {
	return r.db.WithContext(ctx).Model(&models.Departure{}).Create(departure).Error
}
