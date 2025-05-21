package config

import (
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// กำหนดค่า default ถ้าไม่ได้ตั้ง environment
		dsn = "postgresql://postgres:postgres@vote_db:5432/vote_db?sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db = database
}

func DB() *gorm.DB {
	return db
}


func SetupDatabase() {
	// สร้างตาราง Votes ตามโครงสร้าง entity.Votes
	err := db.AutoMigrate(&entity.Votes{})
	if err != nil {
		log.Fatal("Failed to migrate Votes table:", err)
	}

	// ตัวอย่าง seed ข้อมูล (ถ้าต้องการ)
	// ตรวจสอบว่ามีข้อมูลในตารางหรือยัง ถ้ายังให้เพิ่มข้อมูลตัวอย่าง
	var count int64
	db.Model(&entity.Votes{}).Count(&count)
	if count == 0 {
		vote := entity.Votes{
			UserID:      1,
			CandidateID: 1,
			ElectionID:  1,
			Timestamp:   time.Now(),
		}
		if err := db.Create(&vote).Error; err != nil {
			log.Println("Failed to seed Votes:", err)
		}
	}
}