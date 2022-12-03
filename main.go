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

	move := move(moveRequest.You, moveRequest.Board)

	// Respond with the calculated move
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(move); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func move(you Snake, board Board) MoveResponse {
	// Create a list of possible moves
	currentPosition := you.Head
	moves := map[string]Coord{
		"right": {X: currentPosition.X + 1, Y: currentPosition.Y},
		"left":  {X: currentPosition.X - 1, Y: currentPosition.Y},
		"up":    {X: currentPosition.X, Y: currentPosition.Y + 1},
		"down":  {X: currentPosition.X, Y: currentPosition.Y - 1},
	}

	log.Printf("Possible moves: %+v", moves)

	// Filter out moves that would cause the snake to collide with a wall or itself
	for key, move := range moves {
		// Remove moves that would result in colliding with the wall
		if move.X < 0 || move.X >= board.Width || move.Y < 0 || move.Y >= board.Height {
			log.Printf("Removing move: %+v for colliding with wall.", move)
			delete(moves, key)
			continue
		}
		// Remove moves that would result in colliding with any other part of itself
		for _, bodyPart := range you.Body {
			if move.X == bodyPart.X && move.Y == bodyPart.Y {
				log.Printf("Removing move: %+v for colliding with snake.", move)
				delete(moves, key)
			}
		}
	}

	log.Printf("Moves remaining: %+v", moves)

	// When no moves are possible...
	if len(moves) == 0 {
		return MoveResponse{"down", "I HAVE NO MOVES LEFT!!!"}
	}

	directions := make([]string, 0, len(moves))
	for key := range moves {
		directions = append(directions, key)
	}

	// Return a random valid move
	rand.Seed(time.Now().UnixNano())
	return MoveResponse{directions[rand.Intn(len(directions))], ""}
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
