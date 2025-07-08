package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adohong4/carZone/driver"
	carHandler "github.com/adohong4/carZone/handler/car"
	engineHandler "github.com/adohong4/carZone/handler/engine"
	loginHandler "github.com/adohong4/carZone/handler/login"
	middleware "github.com/adohong4/carZone/middleware"
	carService "github.com/adohong4/carZone/service/car"
	engineService "github.com/adohong4/carZone/service/engine"
	carStore "github.com/adohong4/carZone/store/car"
	engineStore "github.com/adohong4/carZone/store/engine"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Cannot install file .env: %v", err)
	}

	traceProvider, err := startTracing()
	if err != nil {
		log.Fatal("Failed to Start Tracing: %v", err)
	}

	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to Shut Down Tracing: %v", err)
		}
	}()

	otel.SetTracerProvider(traceProvider)

	// Connect database
	if err := driver.InitDB(); err != nil {
		log.Fatalf("Unable to initialize the database connection: %v", err)
	}
	defer driver.CloseDB()

	db := driver.GetDB()
	if db == nil {
		log.Fatal("database connection is nil, unable to continue")
	}

	// initialize store, service and handler
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	// initialize router
	router := mux.NewRouter()

	router.Use(otelmux.Middleware("CarZone"))
	router.Use(middleware.MetricMiddleware)

	// excute schema
	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatalf("Cannot excute file schema: %v", err)
	}

	router.HandleFunc("/login", loginHandler.LoginHandler).Methods("POST")

	// Middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Route
	protected.HandleFunc("/cars/{id}", carHandler.GetCarById).Methods("GET")
	protected.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	protected.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")
	protected.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	protected.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	protected.HandleFunc("/engines/{id}", engineHandler.GetEngineByID).Methods("GET")
	protected.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")
	protected.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")
	protected.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	router.Handle("/metrics", promhttp.Handler())

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("localhost run port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func executeSchemaFile(db *sql.DB, fileName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Cannot read file schema %s: %v", fileName, err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("Cannot excute file schema %s: %v", fileName, err)
	}
	log.Printf("Successful fileSchema execution %s", fileName)
	return nil
}

func startTracing() (*trace.TracerProvider, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	expoter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("Error Creating new Exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			expoter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("CarZone"),
			),
		),
	)
	return tracerProvider, nil
}
