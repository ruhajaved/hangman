package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type GameSession struct {
	Word           string
	GuessesLeft    int
	CorrectGuesses []bool
}

var (
	sessions = make(map[string]*GameSession)
	mu       sync.Mutex
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello, world!"})
	})

	r.GET("/word", getWord)
	r.POST("/guess/letter", guessLetter)
	r.POST("/guess/word", guessWord)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getWord(c *gin.Context) {
	// difficulty := c.Query("difficulty")

	// let's use a static word for now
	word := "example"

	// set up the game session, for now let's use a static one
	sessionID := "1"
	mu.Lock()
	sessions[sessionID] = &GameSession{
		Word:           word,
		GuessesLeft:    6,
		CorrectGuesses: make([]bool, len(word)),
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"session_id":  sessionID,
		"word_length": len(word),
	})
}

func guessLetter(c *gin.Context) {
	var request struct {
		SessionID string `json:"session_id"`
		Letter    string `json:"letter"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	session, exists := sessions[request.SessionID]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	letter := request.Letter
	correct := false

	for i, l := range session.Word {
		if string(l) == letter {
			session.CorrectGuesses[i] = true
			correct = true
		}
	}

	if !correct {
		session.GuessesLeft--
	}

	// TODO: Fail if guesses left is 0

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"guesses_left": session.GuessesLeft,
	})
}

func guessWord(c *gin.Context) {
	var request struct {
		SessionID string `json:"session_id"`
		Word      string `json:"word"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	session, exists := sessions[request.SessionID]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	correct := request.Word == session.Word
	if !correct {
		session.GuessesLeft--
	}

	// TODO: Fail if guesses left is 0

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"guesses_left": session.GuessesLeft,
	})
}
