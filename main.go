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
		log.Fatalf("Không thể tải file .env: %v", err)
	}

	// Khởi tạo kết nối cơ sở dữ liệu
	if err := driver.InitDB(); err != nil {
		log.Fatalf("Không thể khởi tạo kết nối cơ sở dữ liệu: %v", err)
	}
	defer driver.CloseDB()

	db := driver.GetDB()
	if db == nil {
		log.Fatal("Kết nối cơ sở dữ liệu là nil, không thể tiếp tục")
	}

	// Khởi tạo store, service và handler
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	// Khởi tạo router
	router := mux.NewRouter()

	// Thực thi schema
	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatalf("Không thể thực thi file schema: %v", err)
	}

	// Định nghĩa các route
	router.HandleFunc("/cars/{id}", carHandler.GetCarById).Methods("GET")
	router.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/engines/{id}", engineHandler.GetEngineByID).Methods("GET")
	router.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")
	router.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")
	router.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	// Lấy port từ biến môi trường
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port mặc định nếu không có trong .env
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Máy chủ đang chạy trên port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func executeSchemaFile(db *sql.DB, fileName string) error {
	if db == nil {
		return fmt.Errorf("kết nối cơ sở dữ liệu là nil")
	}

	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("không thể đọc file schema %s: %v", fileName, err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("không thể thực thi file schema %s: %v", fileName, err)
	}
	log.Printf("Đã thực thi thành công file schema %s", fileName)
	return nil
}
