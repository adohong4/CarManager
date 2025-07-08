package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	log.Println("Đang cố gắng kết nối tới cơ sở dữ liệu với DSN:", connStr)

	// Thử kết nối với cơ sở dữ liệu, retry tối đa 5 lần
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Không thể mở kết nối cơ sở dữ liệu: %v", err)
			return err
		}

		// Kiểm tra kết nối
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Không thể ping cơ sở dữ liệu, thử lại sau 5 giây... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Printf("Không thể kết nối tới cơ sở dữ liệu sau 5 lần thử: %v", err)
		return err
	}

	log.Println("Kết nối cơ sở dữ liệu thành công")
	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		log.Println("Cảnh báo: Kết nối cơ sở dữ liệu chưa được khởi tạo")
	}
	return db
}

func CloseDB() error {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Không thể đóng kết nối cơ sở dữ liệu: %v", err)
			return err
		}
		log.Println("Kết nối cơ sở dữ liệu đã được đóng")
		db = nil
	}
	return nil
}
