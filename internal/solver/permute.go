package solver

func permutations(colors []rune) [][]rune {
	if len(colors) <= 1 {
		return [][]rune{append([]rune(nil), colors...)}
	}
	var out [][]rune
	for i, c := range colors {
		rest := append([]rune(nil), colors[:i]...)
		rest = append(rest, colors[i+1:]...)
		for _, p := range permutations(rest) {
			out = append(out, append([]rune{c}, p...))
		}
	}
	return out
}
