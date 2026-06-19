package puzzle

import (
	"encoding/json"
	"strings"
)

type GridJSON struct {
	Rows   int      `json:"rows"`
	Cols   int      `json:"cols"`
	Grid   []string `json:"grid"`
	Colors []string `json:"colors"`
}

func (b *Board) ToASCII(assign []int) string {
	var sb strings.Builder
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			sb.WriteRune(b.DisplayChar(r, c, assign, b.Colors))
			if c+1 < b.Cols {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (b *Board) ToGridJSON(assign []int) GridJSON {
	rows := make([]string, b.Rows)
	for r := 0; r < b.Rows; r++ {
		var line strings.Builder
		for c := 0; c < b.Cols; c++ {
			line.WriteRune(b.DisplayChar(r, c, assign, b.Colors))
		}
		rows[r] = line.String()
	}
	colors := make([]string, len(b.Colors))
	for i, c := range b.Colors {
		colors[i] = string(c)
	}
	return GridJSON{
		Rows:   b.Rows,
		Cols:   b.Cols,
		Grid:   rows,
		Colors: colors,
	}
}

func (b *Board) MarshalJSON(assign []int) ([]byte, error) {
	return json.Marshal(b.ToGridJSON(assign))
}

func ColorIDs(colors []rune) map[rune]int {
	m := make(map[rune]int, len(colors))
	for i, c := range colors {
		m[c] = i + 1
	}
	return m
}
