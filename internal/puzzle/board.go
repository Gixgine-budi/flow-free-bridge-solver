package puzzle

import "fmt"

type Track int

const (
	TrackNormal Track = iota
	TrackH
	TrackV
)

type Point struct {
	R, C int
}

type CellRef struct {
	R, C  int
	Track Track
}

func (c CellRef) Key() string {
	return fmt.Sprintf("%d,%d,%d", c.R, c.C, c.Track)
}

type Board struct {
	Rows, Cols int
	Wall       [][]bool
	Bridge     [][]bool
	Endpoint   [][]rune
	Colors     []rune
	ColorEnds  map[rune][2]Point
	Slots      []CellRef
	SlotOf     map[string]int
}

func NewBoard(rows, cols int) *Board {
	wall := make([][]bool, rows)
	bridge := make([][]bool, rows)
	ep := make([][]rune, rows)
	for r := 0; r < rows; r++ {
		wall[r] = make([]bool, cols)
		bridge[r] = make([]bool, cols)
		ep[r] = make([]rune, cols)
	}
	return &Board{
		Rows:      rows,
		Cols:      cols,
		Wall:      wall,
		Bridge:    bridge,
		Endpoint:  ep,
		ColorEnds: make(map[rune][2]Point),
		SlotOf:    make(map[string]int),
	}
}

func (b *Board) InBounds(r, c int) bool {
	return r >= 0 && r < b.Rows && c >= 0 && c < b.Cols
}

func (b *Board) IsEndpoint(r, c int) bool {
	return b.Endpoint[r][c] != 0
}

func (b *Board) BuildSlots() {
	b.Slots = nil
	b.SlotOf = make(map[string]int)
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Wall[r][c] {
				continue
			}
			if b.Bridge[r][c] {
				b.addSlot(CellRef{R: r, C: c, Track: TrackH})
				b.addSlot(CellRef{R: r, C: c, Track: TrackV})
			} else {
				b.addSlot(CellRef{R: r, C: c, Track: TrackNormal})
			}
		}
	}
}

func (b *Board) addSlot(ref CellRef) {
	b.SlotOf[ref.Key()] = len(b.Slots)
	b.Slots = append(b.Slots, ref)
}

func (b *Board) Index(ref CellRef) (int, bool) {
	idx, ok := b.SlotOf[ref.Key()]
	return idx, ok
}

func (b *Board) EndpointRef(color rune) (CellRef, CellRef, bool) {
	ends, ok := b.ColorEnds[color]
	if !ok {
		return CellRef{}, CellRef{}, false
	}
	a := CellRef{R: ends[0].R, C: ends[0].C, Track: TrackNormal}
	bb := CellRef{R: ends[1].R, C: ends[1].C, Track: TrackNormal}
	return a, bb, true
}

func (b *Board) Neighbors(ref CellRef) []CellRef {
	var out []CellRef
	dirs := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for _, d := range dirs {
		if n, ok := b.Move(ref, d.R, d.C); ok {
			out = append(out, n)
		}
	}
	return out
}

func (b *Board) Move(from CellRef, dr, dc int) (CellRef, bool) {
	horizontal := dc != 0
	vertical := dr != 0
	if horizontal && vertical {
		return CellRef{}, false
	}

	switch from.Track {
	case TrackH:
		if !horizontal {
			return CellRef{}, false
		}
	case TrackV:
		if !vertical {
			return CellRef{}, false
		}
	}

	nr, nc := from.R+dr, from.C+dc
	if !b.InBounds(nr, nc) || b.Wall[nr][nc] {
		return CellRef{}, false
	}

	if from.Track == TrackH || from.Track == TrackV {
		if b.Bridge[nr][nc] {
			if horizontal {
				return CellRef{R: nr, C: nc, Track: TrackH}, true
			}
			return CellRef{R: nr, C: nc, Track: TrackV}, true
		}
		return CellRef{R: nr, C: nc, Track: TrackNormal}, true
	}

	if b.Bridge[nr][nc] {
		if horizontal {
			return CellRef{R: nr, C: nc, Track: TrackH}, true
		}
		return CellRef{R: nr, C: nc, Track: TrackV}, true
	}
	return CellRef{R: nr, C: nc, Track: TrackNormal}, true
}

func (b *Board) DisplayChar(r, c int, assign []int, colorRunes []rune) rune {
	if b.Wall[r][c] {
		return '#'
	}
	if b.Bridge[r][c] {
		h := CellRef{R: r, C: c, Track: TrackH}
		v := CellRef{R: r, C: c, Track: TrackV}
		hi, hok := b.Index(h)
		vi, vok := b.Index(v)
		var hc, vc rune
		if hok && assign[hi] > 0 {
			hc = colorRunes[assign[hi]-1]
		}
		if vok && assign[vi] > 0 {
			vc = colorRunes[assign[vi]-1]
		}
		if hc != 0 && vc != 0 {
			return '+'
		}
		if hc != 0 {
			return hc
		}
		if vc != 0 {
			return vc
		}
		return 'x'
	}
	ref := CellRef{R: r, C: c, Track: TrackNormal}
	idx, ok := b.Index(ref)
	if !ok {
		return '.'
	}
	if assign[idx] > 0 {
		return colorRunes[assign[idx]-1]
	}
	if b.Endpoint[r][c] != 0 {
		return b.Endpoint[r][c]
	}
	return '.'
}
