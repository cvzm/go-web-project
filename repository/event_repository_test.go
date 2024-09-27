package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cvzm/go-web-project/adapter/storage"
	"github.com/cvzm/go-web-project/domain"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewEventRepository(t *testing.T) {
	gormDB, _ := storage.GetMockDB(t)
	repo := NewEventRepository(gormDB)
	assert.NotNil(t, repo)
	assert.Implements(t, (*domain.EventRepository)(nil), repo)
}

func TestEventRepositorySave(t *testing.T) {
	gormDB, mock := storage.GetMockDB(t)
	repo := NewEventRepository(gormDB)

	createdAt := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "events"`)).
		WithArgs("AWS", "EC2_STARTED", "EC2 instance started",
			pq.StringArray([]string{"A", "B"}),
			createdAt, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	event := &domain.Event{
		Source:            domain.SourceAWS,
		EventType:         "EC2_STARTED",
		Description:       "EC2 instance started",
		AffectedResources: []string{"A", "B"},
		CreatedAt:         createdAt,
	}
	err := repo.Save(event)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), event.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEventRepositoryFindAll(t *testing.T) {
	gormDB, mock := storage.GetMockDB(t)
	repo := NewEventRepository(gormDB)

	timestamp := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "source", "event_type", "description", "created_at", "updated_at"}).
			AddRow(1, "AWS", "EC2_STARTED", "EC2 instance started", timestamp, timestamp).
			AddRow(2, "GCP", "VM_STOPPED", "VM instance stopped", timestamp, timestamp))

	events, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, []domain.Event{
		{
			ID:          1,
			Source:      domain.SourceAWS,
			EventType:   "EC2_STARTED",
			Description: "EC2 instance started",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          2,
			Source:      domain.SourceGCP,
			EventType:   "VM_STOPPED",
			Description: "VM instance stopped",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
	}, events)

	assert.NoError(t, mock.ExpectationsWereMet())
}
