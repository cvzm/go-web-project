package repository

import (
	"github.com/cvzm/go-web-project/domain"

	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Save(event *domain.Event) error {
	return Save(r.db, event)
}

func (r *eventRepository) FindAll() ([]domain.Event, error) {
	return FindAll(r.db, &domain.Event{})
}
