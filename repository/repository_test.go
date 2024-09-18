package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cvzm/go-web-project/adapter/storage"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
}

func TestFindAll(t *testing.T) {
	gormDB, mock := storage.GetMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "A").
			AddRow(2, "B"))

	results, err := FindAll(gormDB, &TestModel{})
	assert.NoError(t, err)
	assert.Equal(t, results, []TestModel{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFind(t *testing.T) {
	gormDB, mock := storage.GetMockDB(t)

	createAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE name = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs("Test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(1, "Test", createAt))

	data, err := Find[TestModel](gormDB, "name", "Test")
	assert.NoError(t, err)
	assert.Equal(t, data, TestModel{
		ID:        1,
		Name:      "Test",
		CreatedAt: createAt,
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave(t *testing.T) {
	gormDB, mock := storage.GetMockDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "test_models"`)).
		WithArgs("A", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
	mock.ExpectCommit()

	event := &TestModel{Name: "A"}
	err := Save(gormDB, event)
	assert.NoError(t, err)
	assert.Equal(t, uint(3), event.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
