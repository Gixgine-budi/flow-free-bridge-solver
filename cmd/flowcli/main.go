package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/solver"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "solve":
		runSolve(os.Args[2:])
	case "bench":
		runBench(os.Args[2:])
	case "verify":
		runVerify(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`flowcli - Flow Free Bridges backtracking solver

Usage:
  flowcli solve [-trace] [-no-connect] <puzzle.txt>
  flowcli bench [-no-connect] <puzzles-dir>
  flowcli verify <puzzle.txt>`)
}

func runSolve(args []string) {
	fs := flag.NewFlagSet("solve", flag.ExitOnError)
	trace := fs.Bool("trace", false, "print backtrack trace")
	noConnect := fs.Bool("no-connect", false, "disable connectivity pruning")
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: flowcli solve [-trace] <puzzle.txt>")
		os.Exit(1)
	}
	board, err := puzzle.ParseFile(fs.Arg(0))
	if err != nil {
		fail(err)
	}
	opts := solver.DefaultOptions()
	opts.Trace = *trace
	if *noConnect {
		opts.PruneConnect = false
	}
	start := time.Now()
	s := solver.New(board, opts)
	ok, assign, stats := s.Solve()
	fmt.Fprintf(os.Stderr, "time=%s\n", time.Since(start))
	if !ok {
		fmt.Println("UNSOLVABLE")
		os.Exit(2)
	}
	fmt.Println(board.ToASCII(assign))
	fmt.Printf("nodes=%d backtracks=%d maxDepth=%d\n", stats.NodesExpanded, stats.Backtracks, stats.MaxDepth)
	if *trace {
		for _, line := range s.Trace {
			fmt.Println(line)
		}
	}
}

func runBench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	noConnect := fs.Bool("no-connect", false, "disable connectivity pruning")
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: flowcli bench <puzzles-dir>")
		os.Exit(1)
	}
	dir := fs.Arg(0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		fail(err)
	}
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"file", "rows", "cols", "colors", "solved", "nodes", "backtracks", "maxDepth", "ms"})
	opts := solver.DefaultOptions()
	if *noConnect {
		opts.PruneConnect = false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		board, err := puzzle.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", e.Name(), err)
			continue
		}
		start := time.Now()
		s := solver.New(board, opts)
		ok, _, stats := s.Solve()
		ms := time.Since(start).Milliseconds()
		solved := "no"
		if ok {
			solved = "yes"
		}
		_ = w.Write([]string{
			e.Name(),
			fmt.Sprint(board.Rows),
			fmt.Sprint(board.Cols),
			fmt.Sprint(len(board.Colors)),
			solved,
			fmt.Sprint(stats.NodesExpanded),
			fmt.Sprint(stats.Backtracks),
			fmt.Sprint(stats.MaxDepth),
			fmt.Sprint(ms),
		})
	}
	w.Flush()
}

func runVerify(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: flowcli verify <puzzle.txt>")
		os.Exit(1)
	}
	board, err := puzzle.ParseFile(args[0])
	if err != nil {
		fail(err)
	}
	opts := solver.DefaultOptions()
	s := solver.New(board, opts)
	ok, assign, stats := s.Solve()
	if !ok {
		fmt.Println("UNSOLVABLE")
		os.Exit(2)
	}
	if err := solver.Verify(board, assign); err != nil {
		fail(err)
	}
	fmt.Println("VALID")
	fmt.Printf("nodes=%d backtracks=%d maxDepth=%d\n", stats.NodesExpanded, stats.Backtracks, stats.MaxDepth)
	fmt.Println(board.ToASCII(assign))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
