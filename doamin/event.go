// Package doamin defines event-related data structures and interfaces

package doamin

import "time"

// EventSource represents the source of an event
type EventSource string

// Constants defining different event sources
const (
	SourceAWS     EventSource = "AWS"
	SourceGCP     EventSource = "GCP"
	SourceAzure   EventSource = "Azure"
	SourceAlibaba EventSource = "Alibaba"
	SourceTencent EventSource = "Tencent"
)

// Event struct defines the properties of an event
type Event struct {
	ID          uint        `gorm:"primaryKey"`
	Source      EventSource `gorm:"type:varchar(255);not null"`
	EventType   string      `gorm:"type:varchar(100);not null"`
	Description string      `gorm:"type:text"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}

// TableName returns the name of the events table
func (Event) TableName() string {
	return "events"
}

// EventRepository defines the interface for event storage
type EventRepository interface {
	Save(event *Event) error
	FindAll() ([]Event, error)
}

// CloudEvent defines the interface for cloud events
type CloudEvent interface {
	Parse() (Event, error)
}

// EventUsecase defines the interface for event use cases
type EventUsecase interface {
	Create(cloudEvent CloudEvent) error
	// GetAllEvents() ([]Event, error)
}

// AWSEvent represents an AWS cloud event
type AWSEvent struct {
	AWSEventID   string    `json:"aws_event_id"`
	AWSEventType string    `json:"aws_event_type"`
	AWSMessage   string    `json:"aws_message"`
	AWSTimestamp time.Time `json:"aws_timestamp"`
}

// Parse implements the CloudEvent interface for AWSEvent
func (a AWSEvent) Parse() (Event, error) {
	return Event{
		Source:      SourceAWS,
		EventType:   a.AWSEventType,
		Description: a.AWSMessage,
		CreatedAt:   a.AWSTimestamp,
	}, nil
}

// GCPEvent represents a GCP cloud event
type GCPEvent struct {
	GCPEventID   string    `json:"gcp_event_id"`
	GCPEventType string    `json:"gcp_event_type"`
	GCPMessage   string    `json:"gcp_message"`
	GCPTimestamp time.Time `json:"gcp_timestamp"`
}

// Parse implements the CloudEvent interface for GCPEvent
func (g GCPEvent) Parse() (Event, error) {
	return Event{
		Source:      SourceGCP,
		EventType:   g.GCPEventType,
		Description: g.GCPMessage,
		CreatedAt:   g.GCPTimestamp,
	}, nil
}
