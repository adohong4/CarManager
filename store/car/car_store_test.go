package car

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adohong4/carZone/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCarById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	carID := uuid.New().String()
	mock.ExpectQuery("SELECT c.id, c.name, c.year, c.brand").
		WithArgs(carID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "year", "brand", "fuel_type", "price", "created_at", "updated_at"}).
			AddRow(carID, "Test Car", 2020, "Test Brand", "Petrol", 20000, time.Now(), time.Now()))

	car, err := store.GetCarById(context.Background(), carID)
	assert.NoError(t, err)
	assert.Equal(t, carID, car.ID)
	assert.Equal(t, "Test Car", car.Name)
}

func TestGetCarByBrand(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	brand := "Test Brand"
	mock.ExpectQuery("SELECT c.id, c.name, c.year, c.brand").
		WithArgs(brand).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "year", "brand", "fuel_type", "price", "created_at", "updated_at"}).
			AddRow(uuid.New().String(), "Test Car 1", 2020, brand, "Petrol", 20000, time.Now(), time.Now()).
			AddRow(uuid.New().String(), "Test Car 2", 2021, brand, "Diesel", 25000, time.Now(), time.Now()))

	cars, err := store.GetCarByBrand(context.Background(), brand, false)
	assert.NoError(t, err)
	assert.Len(t, cars, 2)
}

func TestCreateCar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	engineID := uuid.New()
	carID := uuid.New()
	carReq := &models.CarRequest{
		Name:     "New Car",
		Year:     "2022",
		Brand:    "New Brand",
		FuelType: "Electric",
		Engine: models.Engine{ // Sử dụng giá trị kiểu Engine
			EngineID:      engineID,
			Displacement:  1500, // Thêm các trường cần thiết
			NoOfCylinders: 4,
			CarRange:      300,
		},
		Price: 30000,
	}

	mock.ExpectQuery("SELECT id FROM engine WHERE id = $1").
		WithArgs(carReq.Engine.EngineID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(carReq.Engine.EngineID))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO car").
		WithArgs(carID, carReq.Name, carReq.Year, carReq.Brand, carReq.FuelType, carReq.Engine.EngineID, carReq.Price, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	car, err := store.CreateCar(context.Background(), carReq)
	assert.NoError(t, err)
	assert.Equal(t, carID, car.ID)
}

func TestUpdateCar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	carID := uuid.New().String()
	carReq := &models.CarRequest{
		Name:     "Updated Car",
		Year:     "2023",
		Brand:    "Updated Brand",
		FuelType: "Hybrid",
		Engine: models.Engine{ // Sử dụng giá trị kiểu Engine
			EngineID:      uuid.New(),
			Displacement:  1600,
			NoOfCylinders: 4,
			CarRange:      350,
		},
		Price: 35000,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE car").
		WithArgs(carReq.Name, carReq.Year, carReq.Brand, carReq.FuelType, carReq.Engine.EngineID, carReq.Price, sqlmock.AnyArg(), carID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	car, err := store.UpdateCar(context.Background(), carID, carReq)
	assert.NoError(t, err)
	assert.Equal(t, carID, car.ID)
}

func TestDeleteCar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a database connection", err)
	}
	defer db.Close()

	store := New(db)

	carID := uuid.New().String()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, name, year, brand").
		WithArgs(carID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "year", "brand", "fuel_type", "engine_id", "price", "created_at", "updated_at"}).
			AddRow(carID, "Deleted Car", 2020, "Deleted Brand", "Petrol", uuid.New(), 20000, time.Now(), time.Now()))
	mock.ExpectExec("DELETE FROM car WHERE id = $1").
		WithArgs(carID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	car, err := store.DeleteCar(context.Background(), carID)
	assert.NoError(t, err)
	assert.Equal(t, carID, car.ID)
}
