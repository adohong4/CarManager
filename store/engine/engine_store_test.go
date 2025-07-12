package engine

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adohong4/carZone/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEngineById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	engineID := uuid.New().String()
	mock.ExpectQuery("SELECT id, displacement, no_of_cylinders, car_range").
		WithArgs(engineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "displacement", "no_of_cylinders", "car_range"}).
			AddRow(engineID, 2000, 4, 500))

	engine, err := store.EngineById(context.Background(), engineID)
	assert.NoError(t, err)
	assert.Equal(t, engineID, engine.EngineID.String())
	assert.Equal(t, 2000, engine.Displacement)
}

func TestCreateEngine(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	engineReq := &models.EngineRequest{
		Displacement:  2000,
		NoOfCylinders: 4,
		CarRange:      500,
	}

	engineID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO engine").
		WithArgs(engineID.String(), engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	engine, err := store.CreateEngine(context.Background(), engineReq)
	assert.NoError(t, err)
	assert.Equal(t, engineID.String(), engine.EngineID.String())
	assert.Equal(t, engineReq.Displacement, engine.Displacement)
}

func TestEngineUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	engineID := uuid.New().String()
	engineReq := &models.EngineRequest{
		Displacement:  2500,
		NoOfCylinders: 6,
		CarRange:      600,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE engine").
		WithArgs(engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, engineID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	engine, err := store.EngineUpdate(context.Background(), engineID, engineReq)
	assert.NoError(t, err)
	assert.Equal(t, engineID, engine.EngineID.String())
	assert.Equal(t, engineReq.Displacement, engine.Displacement)
}

func TestEngineDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	engineID := uuid.New().String()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, displacement, no_of_cylinders, car_range").
		WithArgs(engineID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "displacement", "no_of_cylinders", "car_range"}).
			AddRow(engineID, 2000, 4, 500))
	mock.ExpectExec("DELETE FROM engine WHERE id = $1").
		WithArgs(engineID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	engine, err := store.EngineDelete(context.Background(), engineID)
	assert.NoError(t, err)
	assert.Equal(t, engineID, engine.EngineID.String())
}
