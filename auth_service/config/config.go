package config

import (
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// ค่า default ถ้าไม่มี environment variable
		dsn = "postgresql://postgres:postgres@auth_db:5432/auth_db?sslmode=disable"
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
	// สร้างตารางตาม entity.Users
	err := db.AutoMigrate(&entity.Users{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// รหัสผ่าน plain text สำหรับ seed user
	password := "123456"

	// เข้ารหัสรหัสผ่านด้วย bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	user := entity.Users{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@gmail.com",
		Age:       30,
		Password:  string(hashedPassword), // รหัสผ่าน hashed
		Role:      "admin",
		BirthDay:  time.Date(1995, 5, 10, 0, 0, 0, 0, time.UTC),
	}

	// ถ้า user นี้ยังไม่มีใน DB จะสร้างใหม่
	db.FirstOrCreate(&user, entity.Users{Email: "admin@gmail.com"})
}
