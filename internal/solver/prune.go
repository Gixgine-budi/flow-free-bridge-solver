package solver

import (
	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
)

func (s *Solver) connectivityOK(fromColorIdx int) bool {
	if fromColorIdx >= len(s.order) {
		return true
	}

	visited := make(map[string]bool)
	start := s.firstOpenCell()
	if start == (puzzle.CellRef{}) {
		return true
	}

	queue := []puzzle.CellRef{start}
	visited[start.Key()] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		ci, ok := s.Board.Index(cur)
		if !ok || s.Assign[ci] != 0 {
			continue
		}
		for _, n := range s.Board.Neighbors(cur) {
			ni, nok := s.Board.Index(n)
			if !nok || s.Assign[ni] != 0 {
				continue
			}
			if !visited[n.Key()] {
				visited[n.Key()] = true
				queue = append(queue, n)
			}
		}
	}

	for i := fromColorIdx; i < len(s.order); i++ {
		color := s.order[i]
		a, b, ok := s.Board.EndpointRef(color)
		if !ok {
			continue
		}
		inA := visited[a.Key()]
		inB := visited[b.Key()]
		if inA != inB {
			return false
		}
		if inA && inB {
			count := 0
			for j := fromColorIdx; j < len(s.order); j++ {
				c := s.order[j]
				ea, eb, _ := s.Board.EndpointRef(c)
				if visited[ea.Key()] {
					count++
				}
				if visited[eb.Key()] && eb != ea {
					count++
				}
			}
			if count > 2 {
				for j := fromColorIdx; j < len(s.order); j++ {
					c := s.order[j]
					ea, eb, _ := s.Board.EndpointRef(c)
					n := 0
					if visited[ea.Key()] {
						n++
					}
					if visited[eb.Key()] {
						n++
					}
					if n == 1 {
						return false
					}
				}
			}
		}
	}
	return true
}

func (s *Solver) firstOpenCell() puzzle.CellRef {
	for _, slot := range s.Board.Slots {
		idx, ok := s.Board.Index(slot)
		if ok && s.Assign[idx] == 0 {
			return slot
		}
	}
	return puzzle.CellRef{}
}
