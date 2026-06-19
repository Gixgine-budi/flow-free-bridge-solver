package solver

type Stats struct {
	NodesExpanded int `json:"nodesExpanded"`
	Backtracks    int `json:"backtracks"`
	MaxDepth      int `json:"maxDepth"`
}

type Options struct {
	Trace           bool
	PruneOccupancy  bool
	PruneBridgeAxis bool
	PruneConnect    bool
}

func DefaultOptions() Options {
	return Options{
		PruneOccupancy:  true,
		PruneBridgeAxis: true,
		PruneConnect:    true,
	}
}
