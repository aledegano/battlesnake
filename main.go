package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
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
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Food    []Coord `json:"food"`
	Hazards []Coord `json:"hazards"`
	Snakes  []Snake `json:"snakes"`
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

type PossibleMoves map[string]Coord

const CORNER_AVOIDANCE = 2
const FOOD = 1

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

	possibleMoves := possibleMoves(moveRequest.You, moveRequest.Board)

	move := strategy(possibleMoves, moveRequest.You, moveRequest.Board)

	// Respond with the calculated move
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(move); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func possibleMoves(you Snake, board Board) PossibleMoves {
	// Create a list of possible moves
	currentPosition := you.Head
	moves := PossibleMoves{
		"right": {X: currentPosition.X + 1, Y: currentPosition.Y},
		"left":  {X: currentPosition.X - 1, Y: currentPosition.Y},
		"up":    {X: currentPosition.X, Y: currentPosition.Y + 1},
		"down":  {X: currentPosition.X, Y: currentPosition.Y - 1},
	}

	log.Printf("Possible moves: %+v", moves)

	// Filter out moves that would cause the snake to collide with a wall, itself or other snakes
	for key, move := range moves {
		// Remove moves that would result in colliding with the wall
		if move.X < 0 || move.X >= board.Width || move.Y < 0 || move.Y >= board.Height {
			log.Printf("Removing move: %+v for colliding with wall.", move)
			delete(moves, key)
			continue
		}
		// Remove moves that would result in colliding with a hazard
		for iHazard, hazard := range board.Hazards {
			if move.X == hazard.X && move.Y == hazard.Y {
				log.Printf("Removing move: %+v for colliding with hazard %d.", move, iHazard)
				delete(moves, key)
				continue
			}
		}
		// Remove moves that would result in colliding with any snake including itself
		for iSnake, snake := range board.Snakes {
			for _, bodyPart := range snake.Body {
				if move.X == bodyPart.X && move.Y == bodyPart.Y {
					log.Printf("Removing move: %+v for colliding with snake %d.", move, iSnake)
					delete(moves, key)
				}
			}
		}
	}
	log.Printf("Moves remaining: %+v", moves)
	return moves
}

func strategy(moves PossibleMoves, you Snake, board Board) MoveResponse {
	if len(moves) == 0 {
		return MoveResponse{"down", "I HAVE NO MOVES LEFT!!!"}
	}

	movesRating := make(map[string]int, len(moves))
	for key := range moves {
		movesRating[key] = 0
	}

	center := Coord{X: int(board.Width / 2), Y: int(board.Height / 2)}
	// Rate the move that gets the snake closer to the center
	for key, move := range moves {
		if math.Abs(float64(move.X-center.X)) < math.Abs(float64(you.Head.X-center.X)) {
			movesRating[key] += CORNER_AVOIDANCE
			continue
		}
		if math.Abs(float64(move.Y-center.Y)) < math.Abs(float64(you.Head.Y-center.Y)) {
			movesRating[key] += CORNER_AVOIDANCE
		}
	}

	movesRated := make([]string, 0, len(movesRating))
	for move := range movesRating {
		movesRated = append(movesRated, move)
	}
	sort.Slice(movesRated, func(i, j int) bool {
		return movesRating[movesRated[i]] < movesRating[movesRated[j]]
	})
	log.Printf("Rated moves: %+v", movesRated)
	return MoveResponse{movesRated[0],""}
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
