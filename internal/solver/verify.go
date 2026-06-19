package solver

import (
	"fmt"

	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
)

func Verify(b *puzzle.Board, assign []int) error {
	if len(assign) != len(b.Slots) {
		return fmt.Errorf("assignment length mismatch")
	}
	for i, v := range assign {
		if v == 0 {
			return fmt.Errorf("unfilled slot %s", b.Slots[i].Key())
		}
	}

	colorID := puzzle.ColorIDs(b.Colors)
	idColor := make(map[int]rune)
	for c, id := range colorID {
		idColor[id] = c
	}

	for _, color := range b.Colors {
		start, end, ok := b.EndpointRef(color)
		if !ok {
			continue
		}
		si, _ := b.Index(start)
		ei, _ := b.Index(end)
		if assign[si] != colorID[color] || assign[ei] != colorID[color] {
			return fmt.Errorf("color %c endpoints not connected", color)
		}
	}
	return nil
}
