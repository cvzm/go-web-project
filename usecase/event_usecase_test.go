package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/cvzm/go-web-project/doamin"
	domain_mock "github.com/cvzm/go-web-project/doamin/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventUsecase_Save(t *testing.T) {
	mockRepo := new(domain_mock.MockEventRepository)
	usecase := NewEventUsecase(mockRepo)

	t.Run("Successfully save AWS event", func(t *testing.T) {
		awsEvent := doamin.AWSEvent{
			AWSEventID:   "aws-123",
			AWSEventType: "EC2_STARTED",
			AWSMessage:   "EC2 instance started",
			AWSTimestamp: time.Now(),
		}

		expectedEvent := doamin.Event{
			Source:      doamin.SourceAWS,
			EventType:   awsEvent.AWSEventType,
			Description: awsEvent.AWSMessage,
			CreatedAt:   awsEvent.AWSTimestamp,
		}

		mockRepo.On("Save", mock.AnythingOfType("*doamin.Event")).Return(nil).Once()

		err := usecase.Save(awsEvent)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Save", &expectedEvent)
	})

	t.Run("Successfully save GCP event", func(t *testing.T) {
		gcpEvent := doamin.GCPEvent{
			GCPEventID:   "gcp-456",
			GCPEventType: "VM_STOPPED",
			GCPMessage:   "VM instance stopped",
			GCPTimestamp: time.Now(),
		}

		expectedEvent := doamin.Event{
			Source:      doamin.SourceGCP,
			EventType:   gcpEvent.GCPEventType,
			Description: gcpEvent.GCPMessage,
			CreatedAt:   gcpEvent.GCPTimestamp,
		}

		mockRepo.On("Save", mock.AnythingOfType("*doamin.Event")).Return(nil).Once()

		err := usecase.Save(gcpEvent)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Save", &expectedEvent)
	})

	t.Run("Failed to save event", func(t *testing.T) {
		awsEvent := doamin.AWSEvent{
			AWSEventID:   "aws-789",
			AWSEventType: "EC2_TERMINATED",
			AWSMessage:   "EC2 instance terminated",
			AWSTimestamp: time.Now(),
		}

		mockRepo.On("Save", mock.AnythingOfType("*doamin.Event")).Return(errors.New("save failed")).Once()

		err := usecase.Save(awsEvent)

		assert.Error(t, err)
		assert.EqualError(t, err, "save failed")
	})
}
