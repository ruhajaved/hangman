package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type GameSession struct {
	Word           string
	GuessesLeft    int
	CorrectGuesses []bool
}

var (
	sessions = make(map[int]*GameSession)
	mu       sync.Mutex
	conn     *pgx.Conn
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello, world!"})
	})

	r.GET("/start", getWord)
	r.POST("/guess/letter", guessLetter)
	r.POST("/guess/word", guessWord)

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getWord(c *gin.Context) {
	// difficulty := c.Query("difficulty")

	// let's use a static word for now
	word := "example"

	// set up the game session, for now let's use a static one
	var sessionID int
	guessesLeft := 6
	err := conn.QueryRow(context.Background(),
		"INSERT INTO game_sessions (word, guesses_left) VALUES ($1, $2) RETURNING id",
		word, guessesLeft).Scan(&sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game session"})
		return
	}
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
		SessionID int    `json:"session_id"`
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
		SessionID int    `json:"session_id"`
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
