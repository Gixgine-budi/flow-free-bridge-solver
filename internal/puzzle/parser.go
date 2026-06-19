package puzzle

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func ParseFile(path string) (*Board, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(bufio.NewScanner(f))
}

func ParseText(text string) (*Board, error) {
	return Parse(bufio.NewScanner(strings.NewReader(text)))
}

func Parse(sc *bufio.Scanner) (*Board, error) {
	var lines []string
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty puzzle")
	}

	parts := strings.Fields(lines[0])
	if len(parts) < 2 {
		return nil, fmt.Errorf("first line must be: rows cols")
	}
	rows, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid rows: %w", err)
	}
	cols, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid cols: %w", err)
	}
	if rows < 1 || cols < 1 {
		return nil, fmt.Errorf("invalid dimensions")
	}

	b := NewBoard(rows, cols)
	colorCount := make(map[rune]int)

	for _, line := range lines[1:] {
		parts = strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "X", "x":
			if len(parts) != 3 {
				return nil, fmt.Errorf("bridge line: X row col")
			}
			r, c, err := parseCoord(parts[1], parts[2], rows, cols)
			if err != nil {
				return nil, err
			}
			if b.Endpoint[r][c] != 0 {
				return nil, fmt.Errorf("bridge on endpoint at (%d,%d)", r, c)
			}
			b.Bridge[r][c] = true
		case "W", "w":
			if len(parts) != 3 {
				return nil, fmt.Errorf("wall line: W row col")
			}
			r, c, err := parseCoord(parts[1], parts[2], rows, cols)
			if err != nil {
				return nil, err
			}
			b.Wall[r][c] = true
		default:
			if len(parts) != 5 {
				return nil, fmt.Errorf("endpoint line: COLOR r1 c1 r2 c2")
			}
			color := rune(parts[0][0])
			if !unicode.IsLetter(color) {
				return nil, fmt.Errorf("invalid color %q", parts[0])
			}
			r1, c1, err := parseCoord(parts[1], parts[2], rows, cols)
			if err != nil {
				return nil, err
			}
			r2, c2, err := parseCoord(parts[3], parts[4], rows, cols)
			if err != nil {
				return nil, err
			}
			if b.Wall[r1][c1] || b.Wall[r2][c2] {
				return nil, fmt.Errorf("endpoint on wall")
			}
			if b.Bridge[r1][c1] || b.Bridge[r2][c2] {
				return nil, fmt.Errorf("endpoint on bridge")
			}
			if b.Endpoint[r1][c1] != 0 || b.Endpoint[r2][c2] != 0 {
				return nil, fmt.Errorf("overlapping endpoints")
			}
			b.Endpoint[r1][c1] = color
			b.Endpoint[r2][c2] = color
			if colorCount[color] > 0 {
				return nil, fmt.Errorf("duplicate color %c", color)
			}
			b.ColorEnds[color] = [2]Point{{R: r1, C: c1}, {R: r2, C: c2}}
			colorCount[color] = 2
		}
	}

	for color, n := range colorCount {
		if n != 2 {
			return nil, fmt.Errorf("color %c must have exactly two endpoints, got %d", color, n)
		}
	}

	b.Colors = sortedColors(b.ColorEnds)
	b.BuildSlots()
	return b, nil
}

func parseCoord(rs, cs string, maxR, maxC int) (int, int, error) {
	r, err := strconv.Atoi(rs)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row: %w", err)
	}
	c, err := strconv.Atoi(cs)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid col: %w", err)
	}
	if r < 0 || r >= maxR || c < 0 || c >= maxC {
		return 0, 0, fmt.Errorf("coordinate (%d,%d) out of bounds", r, c)
	}
	return r, c, nil
}

func sortedColors(m map[rune][2]Point) []rune {
	seen := make(map[rune]struct{})
	for c := range m {
		seen[c] = struct{}{}
	}
	out := make([]rune, 0, len(seen))
	for c := range seen {
		out = append(out, c)
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
