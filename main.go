package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adohong4/carZone/driver"
	carHandler "github.com/adohong4/carZone/handler/car"
	engineHandler "github.com/adohong4/carZone/handler/engine"
	carService "github.com/adohong4/carZone/service/car"
	engineService "github.com/adohong4/carZone/service/engine"
	carStore "github.com/adohong4/carZone/store/car"
	engineStore "github.com/adohong4/carZone/store/engine"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := mux.NewRouter()

	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatal("Failed to execute schema file:", err)
	}

	router.HandleFunc("/cars/{id}", carHandler.GetCarById).Methods("GET")
	router.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/engines/{id}", engineHandler.GetEngineByID).Methods("GET")
	router.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")
	router.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")
	router.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server is running on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func executeSchemaFile(db *sql.DB, fileName string) error {
	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return err
	}
	return nil
}
