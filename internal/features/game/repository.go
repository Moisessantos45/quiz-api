package game

import (
	"context"
	"fmt"
	"quiz/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) GameRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetNextQuestion(ctx context.Context, gameID uint64, categoryID uint64) (*models.Question, error) {
	var usedIDs []uint64
	r.db.WithContext(ctx).
		Model(&models.GameDetail{}).
		Joins("JOIN departures ON departures.id = game_details.departure_id").
		Where("departures.game_id = ?", gameID).
		Pluck("DISTINCT game_details.question_id", &usedIDs)

	var question models.Question
	query := r.db.WithContext(ctx).
		Preload("Answers").
		Where("category_id = ?", categoryID)

	if len(usedIDs) > 0 {
		query = query.Where("id NOT IN ?", usedIDs)
	}

	if err := query.Order("RANDOM()").First(&question).Error; err != nil {
		return nil, fmt.Errorf("no hay mas preguntas disponibles")
	}
	return &question, nil
}

func (r *PostgresRepository) GetAnswer(ctx context.Context, answerID uint64) (*models.Answer, error) {
	var answer models.Answer
	if err := r.db.WithContext(ctx).Where("id = ?", answerID).First(&answer).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *PostgresRepository) GetDeparture(ctx context.Context, departureID uint64) (*models.Departure, error) {
	var departure models.Departure
	if err := r.db.WithContext(ctx).Where("id = ?", departureID).First(&departure).Error; err != nil {
		return nil, err
	}
	return &departure, nil
}

func (r *PostgresRepository) GetGame(ctx context.Context, gameID uint64) (*models.Game, error) {
	var game models.Game
	if err := r.db.WithContext(ctx).Where("id = ?", gameID).First(&game).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *PostgresRepository) CreateDetail(ctx context.Context, detail *models.GameDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *PostgresRepository) UpdateDepartureScore(ctx context.Context, departureID uint64, score int64, hits int, totalTime int) error {
	return r.db.WithContext(ctx).Model(&models.Departure{}).
		Where("id = ?", departureID).
		Updates(map[string]any{"score": score, "hits": hits, "total_time": totalTime}).Error
}

func (r *PostgresRepository) GetScoreboard(ctx context.Context, gameID uint64) ([]models.Departure, error) {
	var departures []models.Departure
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("game_id = ?", gameID).
		Order("score DESC").
		Find(&departures).Error; err != nil {
		return nil, err
	}
	return departures, nil
}

func (r *PostgresRepository) UpdateGameStatus(ctx context.Context, gameID uint64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Game{}).Where("id = ?", gameID).Update("status", status).Error
}
