package bootstrap

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/cvzm/go-web-project/domain"
)

// SQSConsumer represents a consumer that consumes and processes messages from an AWS SQS queue
type SQSConsumer struct {
	sqsClient    *sqs.Client
	config       *Config
	eventUsecase domain.EventUsecase
}

// NewSQSConsumer creates a new SQSConsumer instance
func NewSQSConsumer(cfg aws.Config, config *Config, eventUsecase domain.EventUsecase) *SQSConsumer {
	return &SQSConsumer{
		sqsClient:    sqs.NewFromConfig(cfg),
		config:       config,
		eventUsecase: eventUsecase,
	}
}

// Start begins the message consumption loop
func (c *SQSConsumer) Start() {
	c.consumeMessages(c.config.SQSQueueURL)
}

// consumeMessages continuously polls the SQS queue and consumes messages
func (c *SQSConsumer) consumeMessages(queueURL string) {
	for {
		messages, err := c.receiveMessages(queueURL)
		if err != nil {
			log.Printf("Error receiving SQS messages: %v", err)
			continue
		}

		for _, message := range messages {
			if err := c.handleMessage(message); err != nil {
				log.Printf("Error handling message: %v", err)
				continue
			}

			if err := c.deleteMessage(queueURL, message); err != nil {
				log.Printf("Error deleting message: %v", err)
			}
		}
	}
}

// receiveMessages retrieves messages from the SQS queue
func (c *SQSConsumer) receiveMessages(queueURL string) ([]types.Message, error) {
	result, err := c.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: c.config.SQSMaxNumberOfMessages,
		WaitTimeSeconds:     c.config.SQSWaitTimeSeconds,
		VisibilityTimeout:   c.config.SQSVisibilityTimeout,
	})
	if err != nil {
		return nil, err
	}
	return result.Messages, nil
}

// handleMessage processes a single SQS message
func (c *SQSConsumer) handleMessage(message types.Message) error {
	var awsEvent domain.AWSEvent
	if err := json.Unmarshal([]byte(*message.Body), &awsEvent); err != nil {
		return err
	}
	return c.eventUsecase.Save(awsEvent)
}

// deleteMessage deletes a processed message from the SQS queue
func (c *SQSConsumer) deleteMessage(queueURL string, message types.Message) error {
	_, err := c.sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}

// initSQSConsumer initializes the SQS consumer and starts it
func initSQSConsumer(cfg *Config, eventUsecase domain.EventUsecase) (*SQSConsumer, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithDefaultRegion(cfg.SQSRegion),
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxBackoffDelay(retry.AddWithMaxAttempts(retry.NewStandard(), 10), 1*time.Minute)
		}),
	)
	if err != nil {
		return nil, err
	}

	consumer := NewSQSConsumer(awsCfg, cfg, eventUsecase)
	return consumer, nil
}
