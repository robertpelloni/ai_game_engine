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

func (g *Grid) GetCellIndex(ix, iy int) int {
	return ix*10000 + iy
}

func (g *Grid) GetCoords(x, y float64) (int, int) {
	return int(x / g.CellSize), int(y / g.CellSize)
}

func (g *Grid) Insert(e Entity, x, y float64) {
	ix, iy := g.GetCoords(x, y)
	idx := g.GetCellIndex(ix, iy)
	g.Cells[idx] = append(g.Cells[idx], e)
}

func (g *Grid) Clear() {
	for k := range g.Cells {
		g.Cells[k] = g.Cells[k][:0]
	}
}

func (g *Grid) GetNearby(x, y float64) []Entity {
	ix, iy := g.GetCoords(x, y)
	var nearby []Entity

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			idx := g.GetCellIndex(ix+dx, iy+dy)
			if entities, ok := g.Cells[idx]; ok {
				nearby = append(nearby, entities...)
			}
		}
	}
	return nearby
}
