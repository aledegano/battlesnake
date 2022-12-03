package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type MoveRequest struct {
	Game  Game  `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

// Game represents a single game of battlesnake.
type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Map     string  `json:"map"`
	Timeout int     `json:"timeout"`
	Source  string  `json:"source"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
	var moveRequest MoveRequest
	if err := decoder.Decode(&moveRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Turn %d: Calculating move for game %s\n", moveRequest.Turn, moveRequest.Game.ID)

	move := move(moveRequest.You.Head, moveRequest.Board)

	// Respond with the calculated move
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(move); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func move(currentPosition Coord, board Board) MoveResponse {
	// Create a list of possible moves
	var moves []map[string]Coord
	moves = append(moves, map[string]Coord{"right": {X: currentPosition.X + 1, Y: currentPosition.Y}})
	moves = append(moves, map[string]Coord{"left": {X: currentPosition.X - 1, Y: currentPosition.Y}})
	moves = append(moves, map[string]Coord{"up": {X: currentPosition.X, Y: currentPosition.Y + 1}})
	moves = append(moves, map[string]Coord{"down": {X: currentPosition.X, Y: currentPosition.Y - 1}})

	// Filter out moves that would cause the snake to collide with a wall
	for i, m := range moves {
		for _, move := range m {
			if move.X < 0 && move.X >= board.Width && move.Y < 0 && move.Y >= board.Height { // Out of bound, pop them from list
				moves = append(moves[:i], moves[i+1:]...)
			}
		}
	}

	// Return a random valid move
	rand.Seed(time.Now().UnixNano())
	for direction := range moves[rand.Intn(len(moves))] {
		return MoveResponse{direction, ""}
	}
	// No valid moves found, default to "down"
	return MoveResponse{"down", ""}
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
