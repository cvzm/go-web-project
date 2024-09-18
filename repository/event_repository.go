package repository

import (
	"github.com/cvzm/go-web-project/doamin"

	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) doamin.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Save(event *doamin.Event) error {
	return Save(r.db, event)
}

func (r *eventRepository) FindAll() ([]doamin.Event, error) {
	return FindAll(r.db, &doamin.Event{})
}
