package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/cvzm/go-web-project/domain"
	domain_mock "github.com/cvzm/go-web-project/domain/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventUsecase_Save(t *testing.T) {
	mockRepo := new(domain_mock.MockEventRepository)
	usecase := NewEventUsecase(mockRepo)

	t.Run("Successfully save AWS event", func(t *testing.T) {
		awsEvent := domain.AWSEvent{
			AWSEventID:   "aws-123",
			AWSEventType: "EC2_STARTED",
			AWSMessage:   "EC2 instance started",
			AWSTimestamp: time.Now(),
		}

		expectedEvent := domain.Event{
			Source:      domain.SourceAWS,
			EventType:   awsEvent.AWSEventType,
			Description: awsEvent.AWSMessage,
			CreatedAt:   awsEvent.AWSTimestamp,
		}

		mockRepo.On("Save", mock.AnythingOfType("*domain.Event")).Return(nil).Once()

		err := usecase.Save(awsEvent)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Save", &expectedEvent)
	})

	t.Run("Successfully save GCP event", func(t *testing.T) {
		gcpEvent := domain.GCPEvent{
			GCPEventID:   "gcp-456",
			GCPEventType: "VM_STOPPED",
			GCPMessage:   "VM instance stopped",
			GCPTimestamp: time.Now(),
		}

		expectedEvent := domain.Event{
			Source:      domain.SourceGCP,
			EventType:   gcpEvent.GCPEventType,
			Description: gcpEvent.GCPMessage,
			CreatedAt:   gcpEvent.GCPTimestamp,
		}

		mockRepo.On("Save", mock.AnythingOfType("*domain.Event")).Return(nil).Once()

		err := usecase.Save(gcpEvent)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Save", &expectedEvent)
	})

	t.Run("Failed to save event", func(t *testing.T) {
		awsEvent := domain.AWSEvent{
			AWSEventID:   "aws-789",
			AWSEventType: "EC2_TERMINATED",
			AWSMessage:   "EC2 instance terminated",
			AWSTimestamp: time.Now(),
		}

		mockRepo.On("Save", mock.AnythingOfType("*domain.Event")).Return(errors.New("save failed")).Once()

		err := usecase.Save(awsEvent)

		assert.Error(t, err)
		assert.EqualError(t, err, "save failed")
	})
}
