package room

import (
	"context"
	"fmt"
	"quiz/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) RoomRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateGame(ctx context.Context, game *models.Game) error {
	return r.db.WithContext(ctx).Create(game).Error
}

func (r *PostgresRepository) GetGameByKey(ctx context.Context, key string) (*models.Game, error) {
	var game models.Game
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&game).Error; err != nil {
		return nil, fmt.Errorf("sala no encontrada")
	}
	return &game, nil
}

func (r *PostgresRepository) GetGameByID(ctx context.Context, id uint64) (*models.Game, error) {
	var game models.Game
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&game).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *PostgresRepository) UpdateGameStatus(ctx context.Context, gameID uint64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Game{}).Where("id = ?", gameID).Update("status", status).Error
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresRepository) CreateDeparture(ctx context.Context, departure *models.Departure) error {
	return r.db.WithContext(ctx).Create(departure).Error
}
