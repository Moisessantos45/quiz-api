package db

import (
	"encoding/json"
	"os"
	"quiz/internal/shared/models"
)

type SeedCategory struct {
	Category  string         `json:"category"`
	Questions []SeedQuestion `json:"questions"`
}

type SeedQuestion struct {
	Text      string       `json:"text"`
	MediaType string       `json:"media_type"`
	Answers   []SeedAnswer `json:"answers"`
}

type SeedAnswer struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

func SeedDatabase() error {
	var countData int64

	if err := DB.Model(&models.Category{}).Count(&countData).Error; err != nil {
		return err
	}

	if countData == 0 {
		if err := DB.Create(&models.Categories).Error; err != nil {
			return err
		}
	}

	var countQuestions int64
	if err := DB.Model(&models.Question{}).Count(&countQuestions).Error; err != nil {
		return err
	}

	if countQuestions == 0 {
		file, err := os.ReadFile("cmd/data/questions.json")
		if err != nil {
			return err
		}

		var seedData []SeedCategory
		if err := json.Unmarshal(file, &seedData); err != nil {
			return err
		}

		for _, sc := range seedData {

			cat := models.Category{}
			if err := DB.Where("UPPER(name) = UPPER(?)", sc.Category).First(&cat).Error; err != nil {
				cat.Name = sc.Category
				if err := DB.Create(&cat).Error; err != nil {
					return err
				}
			}

			for _, sq := range sc.Questions {
				question := models.Question{
					Text:       sq.Text,
					MediaType:  sq.MediaType,
					CategoryID: cat.ID,
				}

				if err := DB.Create(&question).Error; err != nil {
					return err
				}

				for _, sa := range sq.Answers {
					answer := models.Answer{
						Content:    sa.Content,
						IsCorrect:  sa.IsCorrect,
						QuestionID: question.ID,
					}
					if err := DB.Create(&answer).Error; err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
