package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Game represents a single game of battlesnake.
type Game struct {
	ID    string `json:"id"`
	Turn  int    `json:"turn"`
	Board Board  `json:"board"`
	You   Snake  `json:"you"`
}

// Board represents the state of the game board.
type Board struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Food   []Coord `json:"food"`
	Snakes []Snake `json:"snakes"`
}

// Snake represents a single snake in the game.
type Snake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int     `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Tail   Coord   `json:"tail"`
}

// Coord represents a single coordinate on the game board.
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Move represents a move made by a snake.
type Move struct {
	Move string `json:"move"`
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/end", handleEnd)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create the JSON response
	response := map[string]string{
		"apiversion": "1",
		"author":     "aledega",
		"color":      "#3366ff",
		"head":       "caffeine",
		"tail":       "coffee",
		"version":    "0.0.1-beta",
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleStart is the handler for the POST /start endpoint.
func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var game Game
	if err := decoder.Decode(&game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Starting game %s\n", game.ID)
}

// handleMove is the handler for the POST /move endpoint.
func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var game Game
	if err := decoder.Decode(&game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Turn %d: Calculating move for game %s\n", game.Turn, game.ID)

	// Calculate the next move
	var move Move
	if game.You.Head.X == 0 {
		// Snake is against the left wall, move right
		move.Move = "right"
	} else if game.You.Head.X == game.Board.Width-1 {
		// Snake is against the right wall, move left
		move.Move = "left"
	} else if game.You.Head.Y == 0 {
		// Snake is against the top wall, move down
		move.Move = "down"
	} else if game.You.Head.Y == game.Board.Height-1 {
		// Snake is against the bottom wall, move up
		move.Move = "up"
	}

	// Respond with the calculated move
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(move); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleEnd is the handler for the POST /end endpoint.
func handleEnd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var game Game
	if err := decoder.Decode(&game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Ending game %s\n", game.ID)
}
