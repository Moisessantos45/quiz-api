package game

import (
	"context"
	"quiz/internal/shared/models"
)

type GameUseCase struct {
	repo GameRepository
}

func NewGameUseCase(repo GameRepository) GameService {
	return &GameUseCase{repo: repo}
}

func (uc *GameUseCase) GetQuestion(ctx context.Context, gameID uint64) (*models.Question, error) {
	game, err := uc.repo.GetGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetNextQuestion(ctx, gameID, game.CategoryID)
}

func (uc *GameUseCase) SubmitAnswer(ctx context.Context, departureID uint64, questionID uint64, answerID uint64, responseTime int) (*AnswerResult, error) {
	answer, err := uc.repo.GetAnswer(ctx, answerID)
	if err != nil {
		return nil, err
	}

	departure, err := uc.repo.GetDeparture(ctx, departureID)
	if err != nil {
		return nil, err
	}

	detail := &models.GameDetail{
		DepartureID:  departureID,
		QuestionID:   questionID,
		AnswerID:     answerID,
		IsCorrect:    answer.IsCorrect,
		ResponseTime: responseTime,
	}
	if err := uc.repo.CreateDetail(ctx, detail); err != nil {
		return nil, err
	}

	newScore := departure.Score
	newHits := departure.Hits
	if answer.IsCorrect {
		newScore += 10
		newHits++
	}
	newTotalTime := departure.TotalTime + responseTime

	if err := uc.repo.UpdateDepartureScore(ctx, departureID, newScore, newHits, newTotalTime); err != nil {
		return nil, err
	}

	departure.Score = newScore
	departure.Hits = newHits
	departure.TotalTime = newTotalTime

	return &AnswerResult{Detail: detail, Departure: departure}, nil
}

func (uc *GameUseCase) GetScoreboard(ctx context.Context, gameID uint64) ([]models.Departure, error) {
	return uc.repo.GetScoreboard(ctx, gameID)
}

func (uc *GameUseCase) FinishGame(ctx context.Context, gameID uint64) error {
	return uc.repo.UpdateGameStatus(ctx, gameID, "finished")
}
