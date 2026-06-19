package solver

import (
	"testing"

	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
)

func TestTiny2x2(t *testing.T) {
	b, err := puzzle.ParseText(`2 2
R 0 0 0 1
G 1 0 1 1`)
	if err != nil {
		t.Fatal(err)
	}
	ok, assign, _ := New(b, DefaultOptions()).Solve()
	if !ok {
		t.Fatal("expected solvable")
	}
	if err := Verify(b, assign); err != nil {
		t.Fatal(err)
	}
}

func TestClassic3x3(t *testing.T) {
	b, err := puzzle.ParseFile("../../puzzles/classic_3x3.txt")
	if err != nil {
		t.Fatal(err)
	}
	ok, assign, _ := New(b, DefaultOptions()).Solve()
	if !ok {
		t.Fatal("expected solvable")
	}
	if err := Verify(b, assign); err != nil {
		t.Fatal(err)
	}
}

func TestBridge5x5(t *testing.T) {
	b, err := puzzle.ParseFile("../../puzzles/bridge_5x5.txt")
	if err != nil {
		t.Fatal(err)
	}
	ok, assign, stats := New(b, DefaultOptions()).Solve()
	if !ok {
		t.Fatalf("expected solvable, stats=%+v", stats)
	}
	if err := Verify(b, assign); err != nil {
		t.Fatalf("verify: %v\n%s", err, b.ToASCII(assign))
	}
}
