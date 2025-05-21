package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"
	"github.com/gin-gonic/gin"
)

// Candidate struct สำหรับดึงข้อมูลจาก candidate_service
type Candidate struct {
	ID         uint   `json:"ID"`
	Name       string `json:"name"`
	ElectionID uint   `json:"election_id"`
}

// ดึง Elections ทั้งหมด พร้อม candidate_ids (แบบ optimized)
func GetAllElections(c *gin.Context) {
	db := config.DB()
	var elections []entity.Elections
	if err := db.Find(&elections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// เรียก candidate_service แค่ครั้งเดียว
	allCandidates, err := fetchAllCandidates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidates"})
		return
	}

	// map election_id → candidate_ids
	candidateMap := make(map[uint][]uint)
	for _, candidate := range allCandidates {
		candidateMap[candidate.ElectionID] = append(candidateMap[candidate.ElectionID], candidate.ID)
	}

	// ผูก candidate_ids กับแต่ละ election
	for i := range elections {
		elections[i].CandidateIDs = candidateMap[elections[i].ID]
	}

	c.JSON(http.StatusOK, elections)
}

// ดึง election ทีละตัว พร้อม candidate_ids (ยังใช้แบบเก่าได้ เพราะโหลดทีละอัน)
func GetElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	db := config.DB()
	var election entity.Elections
	if err := db.First(&election, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	candidates, err := fetchCandidatesByElectionID(election.ID)
	if err != nil {
		election.CandidateIDs = []uint{}
	} else {
		var candidateIDs []uint
		for _, candidate := range candidates {
			candidateIDs = append(candidateIDs, candidate.ID)
		}
		election.CandidateIDs = candidateIDs
	}

	c.JSON(http.StatusOK, election)
}

func CreateElection(c *gin.Context) {
	var election entity.Elections
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if election.Status == "" {
		election.Status = "pending"
	}
	if election.StartTime.IsZero() {
		election.StartTime = time.Now()
	}

	db := config.DB()
	if err := db.Create(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.ElectionsCreateTotal.Inc()

	var count int64
	db.Model(&entity.Elections{}).Count(&count)
	metrics.ElectionsTotal.Set(float64(count))

	c.JSON(http.StatusCreated, election)
}

func UpdateElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	var election entity.Elections
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var existing entity.Elections
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	election.ID = existing.ID
	if err := db.Save(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.ElectionsUpdateTotal.Inc()

	c.JSON(http.StatusOK, election)
}

func DeleteElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	db := config.DB()
	if err := db.Delete(&entity.Elections{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.ElectionsDeleteTotal.Inc()

	var count int64
	db.Model(&entity.Elections{}).Count(&count)
	metrics.ElectionsTotal.Set(float64(count))

	c.JSON(http.StatusOK, gin.H{"message": "Election deleted"})
}

// ฟังก์ชันใหม่: ดึง candidate ทั้งหมดในครั้งเดียว
func fetchAllCandidates() ([]Candidate, error) {
	url := "http://candidate_service:8003/candidates"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("candidate service returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var candidates []Candidate
	if err := json.Unmarshal(body, &candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}

// ฟังก์ชันเดิม: ดึง candidates ตาม election_id (ยังใช้สำหรับ GetElection)
func fetchCandidatesByElectionID(electionID uint) ([]Candidate, error) {
	url := fmt.Sprintf("http://candidate_service:8003/candidates?election_id=%d", electionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("candidate service returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var candidates []Candidate
	if err := json.Unmarshal(body, &candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}
