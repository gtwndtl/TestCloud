package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// กำหนดค่า default ถ้าไม่ได้ตั้ง environment
		dsn = "postgresql://postgres:postgres@candidate_db:5432/candidate_db?sslmode=disable"
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
	err := db.AutoMigrate(&entity.Candidates{})
	if err != nil {
		log.Fatal("Failed to migrate candidate table:", err)
	}

	candidates := []entity.Candidates{
		{Name: "Alice Smith", ElectionID: 1},
		{Name: "Bob Johnson", ElectionID: 1},
		{Name: "Charlie Brown", ElectionID: 1},
		{Name: "Diana Clark", ElectionID: 2},
		{Name: "Ethan Wilson", ElectionID: 2},
		{Name: "Fiona Davis", ElectionID: 2},
		{Name: "George Martin", ElectionID: 3},
		{Name: "Hannah Moore", ElectionID: 3},
		{Name: "Ian Taylor", ElectionID: 3},
		{Name: "Julia Anderson", ElectionID: 4},
		{Name: "Kevin Thomas", ElectionID: 4},
		{Name: "Laura White", ElectionID: 4},
		{Name: "Michael Harris", ElectionID: 5},
		{Name: "Nina Lewis", ElectionID: 5},
		{Name: "Oscar Young", ElectionID: 5},
		{Name: "Paula King", ElectionID: 6},
		{Name: "Quentin Scott", ElectionID: 6},
		{Name: "Rachel Green", ElectionID: 6},
		{Name: "Steve Hall", ElectionID: 7},
		{Name: "Tina Adams", ElectionID: 7},
		{Name: "Umar Baker", ElectionID: 7},
		{Name: "Vera Mitchell", ElectionID: 8},
		{Name: "William Carter", ElectionID: 8},
		{Name: "Xena Roberts", ElectionID: 8},
		{Name: "Yusuf Phillips", ElectionID: 9},
		{Name: "Thuwanon", ElectionID: 9},
		{Name: "Ryan", ElectionID: 9},
	}

	for _, c := range candidates {
		db.FirstOrCreate(&c, entity.Candidates{Name: c.Name})
	}
}
