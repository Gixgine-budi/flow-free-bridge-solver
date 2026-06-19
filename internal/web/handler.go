package web

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/solver"
)

//go:embed static/*
var staticFS embed.FS

type solveRequest struct {
	Puzzle string `json:"puzzle"`
}

type solveResponse struct {
	Solved     bool            `json:"solved"`
	Grid       puzzle.GridJSON `json:"grid,omitempty"`
	Stats      solver.Stats    `json:"stats"`
	TimeMS     int64           `json:"timeMs"`
	Error      string          `json:"error,omitempty"`
	InputGrid  puzzle.GridJSON `json:"inputGrid"`
}

func Handler() http.Handler {
	mux := http.NewServeMux()
	sub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/solve", handleSolve)
	return mux
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req solveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, solveResponse{Error: "invalid JSON body"})
		return
	}
	board, err := puzzle.ParseText(req.Puzzle)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, solveResponse{Error: err.Error()})
		return
	}
	empty := make([]int, len(board.Slots))
	inputGrid := board.ToGridJSON(empty)

	start := time.Now()
	s := solver.New(board, solver.DefaultOptions())
	ok, assign, stats := s.Solve()
	elapsed := time.Since(start).Milliseconds()

	resp := solveResponse{
		Solved:    ok,
		Stats:     stats,
		TimeMS:    elapsed,
		InputGrid: inputGrid,
	}
	if ok {
		resp.Grid = board.ToGridJSON(assign)
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
