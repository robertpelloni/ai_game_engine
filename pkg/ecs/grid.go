package ecs

type Grid struct {
	CellSize float64
	Cells    map[int][]Entity
}

func NewGrid(cellSize float64) *Grid {
	return &Grid{
		CellSize: cellSize,
		Cells:    make(map[int][]Entity),
	}
}

func (g *Grid) GetCellIndex(x, y float64) int {
	// Simple 2D to 1D mapping
	ix := int(x / g.CellSize)
	iy := int(y / g.CellSize)
	return ix*10000 + iy
}

func (g *Grid) Insert(e Entity, x, y float64) {
	idx := g.GetCellIndex(x, y)
	g.Cells[idx] = append(g.Cells[idx], e)
}

func (g *Grid) Clear() {
	for k := range g.Cells {
		g.Cells[k] = g.Cells[k][:0]
	}
}

func (g *Grid) GetNearby(x, y float64) []Entity {
	idx := g.GetCellIndex(x, y)
	// For simplicity, just return the current cell.
	// A better version would check 8 neighbors.
	return g.Cells[idx]
}
