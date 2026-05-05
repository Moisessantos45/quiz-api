package game

import (
	"context"
	"quiz/internal/shared/models"
)

type AnswerResult struct {
	Detail    *models.GameDetail `json:"detail"`
	Departure *models.Departure  `json:"departure"`
}

type GameRepository interface {
	GetNextQuestion(ctx context.Context, gameID uint64, categoryID uint64) (*models.Question, error)
	GetAnswer(ctx context.Context, answerID uint64) (*models.Answer, error)
	GetDeparture(ctx context.Context, departureID uint64) (*models.Departure, error)
	GetGame(ctx context.Context, gameID uint64) (*models.Game, error)
	CreateDetail(ctx context.Context, detail *models.GameDetail) error
	UpdateDepartureScore(ctx context.Context, departureID uint64, score int64, hits int, totalTime int) error
	GetScoreboard(ctx context.Context, gameID uint64) ([]models.Departure, error)
	UpdateGameStatus(ctx context.Context, gameID uint64, status string) error
}

type GameService interface {
	GetQuestion(ctx context.Context, gameID uint64) (*models.Question, error)
	SubmitAnswer(ctx context.Context, departureID uint64, questionID uint64, answerID uint64, responseTime int) (*AnswerResult, error)
	GetScoreboard(ctx context.Context, gameID uint64) ([]models.Departure, error)
	FinishGame(ctx context.Context, gameID uint64) error
}
