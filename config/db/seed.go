package db

import "quiz/internal/shared/models"

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

	// esta cmando proque aun no hay que crearl
	// if err := DB.Model(&models.Option{}).Count(&countData).Error; err != nil {
	// 	return err
	// }

	return nil
}
