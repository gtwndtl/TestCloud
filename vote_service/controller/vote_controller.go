package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"
)

// Structs สำหรับ response จากแต่ละ service
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Candidate struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Election struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func fetchUser(userID uint) (*User, error) {
	url := fmt.Sprintf("http://user_service:8001/user/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch user: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("user service returned status %d", resp.StatusCode)
		log.Println(errMsg)
		return nil, errMsg
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode user response: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func fetchCandidate(candidateID uint) (*Candidate, error) {
	url := fmt.Sprintf("http://candidate_service:8003/candidate/%d", candidateID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch candidate: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("candidate service returned status %d", resp.StatusCode)
		log.Println(errMsg)
		return nil, errMsg
	}

	var candidate Candidate
	if err := json.NewDecoder(resp.Body).Decode(&candidate); err != nil {
		log.Printf("Failed to decode candidate response: %v\n", err)
		return nil, err
	}

	return &candidate, nil
}

func fetchElection(electionID uint) (*Election, error) {
	url := fmt.Sprintf("http://election_service:8002/election/%d", electionID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch election: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("election service returned status %d", resp.StatusCode)
		log.Println(errMsg)
		return nil, errMsg
	}

	var election Election
	if err := json.NewDecoder(resp.Body).Decode(&election); err != nil {
		log.Printf("Failed to decode election response: %v\n", err)
		return nil, err
	}

	return &election, nil
}

func GetVoteWithDetails(c *gin.Context) {
	db := config.DB()
	var votes []entity.Votes

	if err := db.Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type VoteWithDetails struct {
		entity.Votes
		User      *User      `json:"user,omitempty"`
		Candidate *Candidate `json:"candidate,omitempty"`
		Election  *Election  `json:"election,omitempty"`
	}

	var result []VoteWithDetails
	for _, v := range votes {
		user, errUser := fetchUser(v.UserID)
		if errUser != nil {
			log.Printf("Warning: cannot fetch user %d: %v", v.UserID, errUser)
		}

		candidate, errCandidate := fetchCandidate(v.CandidateID)
		if errCandidate != nil {
			log.Printf("Warning: cannot fetch candidate %d: %v", v.CandidateID, errCandidate)
		}

		election, errElection := fetchElection(v.ElectionID)
		if errElection != nil {
			log.Printf("Warning: cannot fetch election %d: %v", v.ElectionID, errElection)
		}

		result = append(result, VoteWithDetails{
			Votes:     v,
			User:      user,
			Candidate: candidate,
			Election:  election,
		})
	}

	c.JSON(http.StatusOK, result)
}

func CreateVote(c *gin.Context) {
	var vote entity.Votes

	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vote.Timestamp = time.Now()

	db := config.DB()
	if err := db.Create(&vote).Error; err != nil {
		log.Printf("Failed to create vote: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	if metrics.VotesCreatedTotal != nil {
		metrics.VotesCreatedTotal.Inc()
	}

	c.JSON(http.StatusCreated, vote)
}

func GetAllVotes(c *gin.Context) {
	db := config.DB()
	var votes []entity.Votes

	if err := db.Find(&votes).Error; err != nil {
		log.Printf("Failed to get votes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

func GetVotesByCandidate(c *gin.Context) {
	candidateIDStr := c.Param("candidate_id")
	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	db := config.DB()
	var votes []entity.Votes

	if err := db.Where("candidate_id = ?", candidateID).Find(&votes).Error; err != nil {
		log.Printf("Failed to get votes by candidate: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

func DeleteVote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote ID"})
		return
	}

	db := config.DB()
	if err := db.Delete(&entity.Votes{}, id).Error; err != nil {
		log.Printf("Failed to delete vote: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
		return
	}

	if metrics.VotesDeletedTotal != nil {
		metrics.VotesDeletedTotal.Inc()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
}
