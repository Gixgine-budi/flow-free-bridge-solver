package solver

import (
	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/puzzle"
)

type Solver struct {
	Board   *puzzle.Board
	Assign  []int
	ColorID map[rune]int
	Stats   Stats
	Opts    Options
	Trace   []string
	order   []rune
}

func New(b *puzzle.Board, opts Options) *Solver {
	colorID := puzzle.ColorIDs(b.Colors)
	assign := make([]int, len(b.Slots))
	s := &Solver{
		Board:   b,
		Assign:  assign,
		ColorID: colorID,
		Opts:    opts,
	}
	s.order = s.orderColors()
	return s
}

func (s *Solver) Solve() (bool, []int, Stats) {
	s.Stats = Stats{}
	s.Trace = nil
	perms := permutations(s.Board.Colors)
	for _, order := range perms {
		s.order = order
		for i := range s.Assign {
			s.Assign[i] = 0
		}
		s.Trace = nil
		localStats := Stats{}
		s.Stats = localStats
		if s.solveColor(0) {
			return true, append([]int(nil), s.Assign...), s.Stats
		}
	}
	return false, nil, s.Stats
}

func (s *Solver) solveColor(colorIdx int) bool {
	if colorIdx == len(s.order) {
		return s.allFilled()
	}
	color := s.order[colorIdx]
	start, end, ok := s.Board.EndpointRef(color)
	if !ok {
		return false
	}
	s.clearColor(color)
	return s.walk(color, start, end, start, puzzle.CellRef{}, colorIdx, 1)
}

func (s *Solver) clearColor(color rune) {
	id := s.ColorID[color]
	for i := range s.Assign {
		if s.Assign[i] == id {
			s.Assign[i] = 0
		}
	}
}

func (s *Solver) allFilled() bool {
	for _, v := range s.Assign {
		if v == 0 {
			return false
		}
	}
	return true
}

func (s *Solver) emptyCount() int {
	n := 0
	for _, v := range s.Assign {
		if v == 0 {
			n++
		}
	}
	return n
}

func (s *Solver) walk(color rune, start, end, cur, prev puzzle.CellRef, colorIdx, depth int) bool {
	s.Stats.NodesExpanded++
	if depth > s.Stats.MaxDepth {
		s.Stats.MaxDepth = depth
	}

	idx, ok := s.Board.Index(cur)
	if !ok {
		return false
	}
	if s.Assign[idx] != 0 {
		return false
	}

	isLast := colorIdx == len(s.order)-1

	s.Assign[idx] = s.ColorID[color]
	if s.Opts.Trace {
		s.Trace = append(s.Trace, formatStep(color, cur, depth))
	}

	if cur == end {
		if isLast {
			if !s.allFilled() {
				s.Assign[idx] = 0
				s.Stats.Backtracks++
				return false
			}
			return true
		}
		if s.Opts.PruneConnect && !s.connectivityOK(colorIdx+1) {
			s.Assign[idx] = 0
			s.Stats.Backtracks++
			return false
		}
		if s.solveColor(colorIdx + 1) {
			return true
		}
		s.Assign[idx] = 0
		s.Stats.Backtracks++
		return false
	}

	for _, n := range s.sortedNeighbors(cur, end) {
		if n == prev {
			continue
		}
		if isLast && n == end && s.emptyCount() > 1 {
			continue
		}
		if !s.canMove(color, n) {
			continue
		}
		if s.walk(color, start, end, n, cur, colorIdx, depth+1) {
			return true
		}
	}

	s.Assign[idx] = 0
	s.Stats.Backtracks++
	return false
}

func (s *Solver) orderColors() []rune {
	out := append([]rune(nil), s.Board.Colors...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if s.colorDistance(out[j]) < s.colorDistance(out[i]) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func (s *Solver) colorDistance(color rune) int {
	ends := s.Board.ColorEnds[color]
	dr := ends[0].R - ends[1].R
	dc := ends[0].C - ends[1].C
	if dr < 0 {
		dr = -dr
	}
	if dc < 0 {
		dc = -dc
	}
	return dr + dc
}

func (s *Solver) sortedNeighbors(cur, target puzzle.CellRef) []puzzle.CellRef {
	ns := s.Board.Neighbors(cur)
	for i := 0; i < len(ns); i++ {
		for j := i + 1; j < len(ns); j++ {
			if manhattan(ns[j], target) < manhattan(ns[i], target) {
				ns[i], ns[j] = ns[j], ns[i]
			}
		}
	}
	return ns
}

func manhattan(a, b puzzle.CellRef) int {
	dr := a.R - b.R
	dc := a.C - b.C
	if dr < 0 {
		dr = -dr
	}
	if dc < 0 {
		dc = -dc
	}
	return dr + dc
}

func (s *Solver) canMove(color rune, to puzzle.CellRef) bool {
	toIdx, ok := s.Board.Index(to)
	if !ok {
		return false
	}
	if s.Assign[toIdx] != 0 {
		return false
	}
	if to.Track == puzzle.TrackNormal {
		if ep := s.Board.Endpoint[to.R][to.C]; ep != 0 && ep != color {
			return false
		}
	}
	return true
}

func formatStep(color rune, cur puzzle.CellRef, depth int) string {
	return string(color) + "@" + cur.Key() + " d=" + itoa(depth)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
