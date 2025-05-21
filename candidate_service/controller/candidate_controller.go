package controller

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "example.com/se/config"
    "example.com/se/entity"
    "example.com/se/metrics"
)

// Struct สำหรับส่งกลับ Candidate พร้อมข้อมูล Election
type CandidateWithElection struct {
    ID         uint `json:"id"`
    Name       string `json:"name"`
    ElectionID uint `json:"election_id"`
}

// GET /candidates
func GetAllCandidates(c *gin.Context) {
    db := config.DB()
    var candidates []entity.Candidates

    if err := db.Find(&candidates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }

    metrics.CandidatesTotal.Set(float64(len(candidates)))

    var result []CandidateWithElection
    for _, cand := range candidates {
        result = append(result, CandidateWithElection{
            ID:         cand.ID,
            Name:       cand.Name,
            ElectionID: cand.ElectionID,
        })
    }

    c.JSON(http.StatusOK, result)
}

// POST /candidates
func CreateCandidate(c *gin.Context) {
    db := config.DB()

    var candidate entity.Candidates
    if err := c.ShouldBindJSON(&candidate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := db.Create(&candidate).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    metrics.CandidatesCreateTotal.Inc()
    var count int64
    db.Model(&entity.Candidates{}).Count(&count)
    metrics.CandidatesTotal.Set(float64(count))

    c.JSON(http.StatusCreated, candidate)
}

// PUT /candidates/:id
func UpdateCandidate(c *gin.Context) {
    db := config.DB()
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
        return
    }

    var updated entity.Candidates
    if err := c.ShouldBindJSON(&updated); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existing entity.Candidates
    if err := db.First(&existing, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
        return
    }

    updated.ID = existing.ID
    if err := db.Save(&updated).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    metrics.CandidatesUpdateTotal.Inc()
    c.JSON(http.StatusOK, updated)
}

// DELETE /candidates/:id
func DeleteCandidate(c *gin.Context) {
    db := config.DB()
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
        return
    }

    if err := db.Delete(&entity.Candidates{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    metrics.CandidatesDeleteTotal.Inc()

    var count int64
    db.Model(&entity.Candidates{}).Count(&count)
    metrics.CandidatesTotal.Set(float64(count))

    c.JSON(http.StatusOK, gin.H{"message": "Candidate deleted successfully"})
}

// GET /candidates/:id
func GetCandidate(c *gin.Context) {
    db := config.DB()
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
        return
    }

    var candidate entity.Candidates
    if err := db.First(&candidate, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
        return
    }

    response := CandidateWithElection{
        ID:         candidate.ID,
        Name:       candidate.Name,
        ElectionID: candidate.ElectionID,
    }

    c.JSON(http.StatusOK, response)
}