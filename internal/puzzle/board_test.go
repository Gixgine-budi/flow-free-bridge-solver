package puzzle

import "testing"

func TestBridgeMoveVertical(t *testing.T) {
	b, err := ParseText(`3 3
R 0 1 2 1
G 1 0 1 2
X 1 1`)
	if err != nil {
		t.Fatal(err)
	}
	cur := CellRef{R: 0, C: 1, Track: TrackNormal}
	v, ok := b.Move(cur, 1, 0)
	if !ok || v.Track != TrackV || v.R != 1 || v.C != 1 {
		t.Fatalf("expected bridge V at (1,1), got %+v ok=%v", v, ok)
	}
	end, ok := b.Move(v, 1, 0)
	if !ok || end != (CellRef{R: 2, C: 1, Track: TrackNormal}) {
		t.Fatalf("expected (2,1) normal, got %+v ok=%v", end, ok)
	}
}

func TestBridgeMoveHorizontal(t *testing.T) {
	b, err := ParseText(`3 3
R 0 1 2 1
G 1 0 1 2
X 1 1`)
	if err != nil {
		t.Fatal(err)
	}
	cur := CellRef{R: 1, C: 0, Track: TrackNormal}
	h, ok := b.Move(cur, 0, 1)
	if !ok || h.Track != TrackH {
		t.Fatalf("expected bridge H, got %+v", h)
	}
}
