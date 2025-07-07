package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	Brand    string  `json:"brand"`
	FuelType string  `json:"fuel_type"`
	Engine   Engine  `json:"engine"`
	Price    float64 `json:"price"`
}

func validateRequest(carRequest CarRequest) error {
	if err := validateName(carRequest.Name); err != nil {
		return err
	}
	if err := validateYear(carRequest.Year); err != nil {
		return err
	}
	if err := validatedBrand(carRequest.Brand); err != nil {
		return err
	}
	if err := ValidateFuelType(carRequest.FuelType); err != nil {
		return err
	}
	if err := valdateEngine(carRequest.Engine); err != nil {
		return err
	}
	if err := validatePrice(carRequest.Price); err != nil {
		return err
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("Name is Required")
	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("Year is Required")
	}

	_, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("Year must be a valid number")
	}
	currentYear := time.Now().Year()
	yearInt, _ := strconv.Atoi(year)
	if yearInt < 1886 || yearInt > currentYear {
		return errors.New("Year must be between 1886 and the current year")
	}
	return nil
}

func validatedBrand(brand string) error {
	if brand == "" {
		return errors.New("Brand is Required")
	}
	return nil
}

func ValidateFuelType(fuelType string) error {
	validFuelTypes := []string{"Petrol", "Diesel", "Electric", "Hybrid"}
	for _, validFuelType := range validFuelTypes {
		if fuelType == validFuelType {
			return nil
		}
	}
	return errors.New("FeulType must be one of Petrol, Diesel, Electric, or Hybrid")
}

func valdateEngine(engine Engine) error {
	if engine.EngineID == uuid.Nil {
		return errors.New("Engine ID is Required")
	}
	if engine.Displacement <= 0 {
		return errors.New("Displacement must be greater than 0")
	}
	if engine.NoOfCylinders <= 0 {
		return errors.New("NoOfCylinders must be greater than 0")
	}
	if engine.CarRange <= 0 {
		return errors.New("CarRange must be greater than 0")
	}
	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("Price must be greater than 0")
	}
	return nil
}
