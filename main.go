package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const MAX_GUESSES = 6

type GameSession struct {
	Word             string
	IncorrectGuesses int
	WordFilled       []byte
}

var session GameSession


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
	word := "example"
	session = GameSession{
		Word:             word,
		IncorrectGuesses: 0,
		WordFilled:       byteSlice(word),
	}

	c.JSON(http.StatusOK, gin.H{"word": word})
}

func byteSlice(word string) []byte {
	length := len(word)
	data := make([]byte, length)
	for i := range data {
		data[i] = '*'
	}
	return data
}

func guessLetter(c *gin.Context) {
	var request struct {
		Letter string `json:"letter"`
	}
    err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter := request.Letter
	correct := false

	for i, l := range session.Word {
		if string(l) == letter {
			session.WordFilled[i] = byte(l)
			correct = true
		}
	}

	if !correct {
		session.IncorrectGuesses++
	}
    if string(session.WordFilled) == session.Word {
        c.JSON(http.StatusOK, gin.H{
            "message": "You Won!",
        })
        return
    }

	if session.IncorrectGuesses == MAX_GUESSES {
        c.JSON(http.StatusOK, gin.H{
            "message": "You Lost! Start Again.",
            "current_state": string(session.WordFilled),
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"guesses_left": MAX_GUESSES - session.IncorrectGuesses,
        "current_state": string(session.WordFilled),
	})
}

func guessWord(c *gin.Context) {
	var request struct {
		Word string `json:"word"`
	}
    err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word := request.Word
	correct := false

    if word == session.Word {
        c.JSON(http.StatusOK, gin.H{
            "message": "You Won!",
        })
        return
    }

    session.IncorrectGuesses++;

	if session.IncorrectGuesses == MAX_GUESSES {
        c.JSON(http.StatusOK, gin.H{
            "message": "You Lost! Start Again.",
            "current_state": string(session.WordFilled),
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"guesses_left": MAX_GUESSES - session.IncorrectGuesses,
        "current_state": string(session.WordFilled),
	})
}
