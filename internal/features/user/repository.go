package user

import (
	"quiz/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User

	if err := r.db.Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Create(user *models.User) error {
	return r.db.Model(&models.User{}).Create(user).Error
}

func (r *PostgresRepository) CreateDeparture(departure *models.Departure) error {
	return r.db.Model(&models.Departure{}).Create(departure).Error
}
