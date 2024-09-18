package usecase

import (
	"github.com/cvzm/go-web-project/doamin"
	"github.com/stretchr/testify/mock"
)

// MockEventUsecase is a mock implementation of doamin.EventUsecase
type MockEventUsecase struct {
	mock.Mock
}

func (m *MockEventUsecase) Create(cloudEvent doamin.CloudEvent) error {
	args := m.Called(cloudEvent)
	return args.Error(0)
}
