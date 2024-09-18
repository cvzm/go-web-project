package repository

import (
	"github.com/cvzm/go-web-project/doamin"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	mock.Mock
}

// Save mocks the method for saving an event
func (m *MockEventRepository) Save(event *doamin.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// FindAll mocks the method for finding all events
func (m *MockEventRepository) FindAll() ([]doamin.Event, error) {
	args := m.Called()
	return args.Get(0).([]doamin.Event), args.Error(1)
}
