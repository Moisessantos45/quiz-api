package category

import (
	"context"
	"quiz/internal/shared/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
}

type CategoryService interface {
	GetAll(ctx context.Context) ([]models.Category, error)
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) CategoryRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

type CategoryUseCase struct {
	repo CategoryRepository
}

func NewCategoryUseCase(repo CategoryRepository) CategoryService {
	return &CategoryUseCase{repo: repo}
}

func (uc *CategoryUseCase) GetAll(ctx context.Context) ([]models.Category, error) {
	return uc.repo.GetAll(ctx)
}
