package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=election_db user=postgres password=postgres dbname=election_db port=5432 sslmode=disable TimeZone=Asia/Bangkok"
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
	err := db.AutoMigrate(&entity.Elections{})
	if err != nil {
		log.Fatal("Failed to migrate Elections table:", err)
	}

	elections := []entity.Elections{
		{
			Title:       "University Student Council 2025",
			Description: "เลือกตั้งสภานักศึกษา ประจำปีการศึกษา 2568",
			StartTime:   time.Now().Add(24 * time.Hour),
			EndTime:     time.Now().Add(48 * time.Hour),
			Status:      "upcoming",
			CandidateIDs: []uint{1, 2, 3}, // ถ้า struct ไม่มีให้ลบหรือแก้ไขตามจริง
		},
		{
			Title:       "เลือกตั้งประธานชมรมวิทยาศาสตร์",
			Description: "คัดเลือกหัวหน้าชมรมวิทยาศาสตร์ปี 2568",
			StartTime:   time.Now().Add(24 * time.Hour),
			EndTime:     time.Now().Add(48 * time.Hour),
			Status:      "upcoming",
			CandidateIDs: []uint{3, 4, 5},
		},
		{
			Title:       "เลือกตั้งหัวหน้าแผนกไอที บริษัท TechNova",
			Description: "โหวตเลือกหัวหน้าแผนกจากพนักงานทั้งหมด",
			StartTime:   time.Now().Add(24 * time.Hour),
			EndTime:     time.Now().Add(48 * time.Hour),
			Status:      "upcoming",
			CandidateIDs: []uint{6, 7, 8},
		},
		{
			Title:       "เลือกตั้งตัวแทนชั้นปีวิศวกรรมศาสตร์",
			Description: "เลือกตัวแทนชั้นปีปี 3 ภาควิชาวิศวกรรมคอมพิวเตอร์",
			StartTime:   time.Now().Add(24 * time.Hour),
			EndTime:     time.Now().Add(48 * time.Hour),
			Status:      "upcoming",
			CandidateIDs: []uint{9, 10, 11},
		},
	}

	for _, elec := range elections {
		db.FirstOrCreate(&elec, entity.Elections{Title: elec.Title})
	}
}
