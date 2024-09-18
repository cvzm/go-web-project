package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// FindAll retrieves all records from the database that match the provided parameter.
// It returns a slice of the generic type T and an error.
func FindAll[T any](db *gorm.DB, param *T) ([]T, error) {
	data := []T{}
	err := db.Where(param).Find(&data).Error
	return data, err
}

// Find retrieves a record by its field.
func Find[T any](db *gorm.DB, field string, value any) (T, error) {
	var data T
	err := db.First(&data, fmt.Sprintf("%s = ?", field), value).Error
	if err == gorm.ErrRecordNotFound {
		return data, nil
	}
	return data, err
}

// Save persists one or more records.
func Save[T any](db *gorm.DB, data ...*T) error {
	if len(data) == 0 {
		return nil
	}
	total := len(data)
	batchSize := db.Config.CreateBatchSize
	if batchSize == 0 || total <= batchSize {
		return db.Save(data).Error
	}

	// The save creation is not batch processed, so the active batch call
	tx := db.Begin()

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		if err := tx.Save(data[i:end]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
