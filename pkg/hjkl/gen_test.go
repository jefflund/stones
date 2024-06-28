package hjkl

import "testing"

func TestGenTileGrid_GenTile(t *testing.T) {
	numCalls := 0
	offsets := make(map[Vector]struct{})
	GenTileGrid(20, 10, func(o Vector) *Tile {
		numCalls++
		offsets[o] = struct{}{}
		return NewTile(o)
	})
	if numCalls != 20*10 {
		t.Fatal("GenTile called incorrect number of times")
	}
	for x := 0; x < 20; x++ {
		for y := 0; y < 10; y++ {
			if _, ok := offsets[Vector{x, y}]; !ok {
				t.Fatal("GenTile missing an offset")
			}
		}
	}
}

func TestGenTileGrid_Offsets(t *testing.T) {
	grid := GenTileGrid(20, 10, NewTile)
	if len(grid) != 20*10 {
		t.Fatal("GetTileGrid created incorrect grid size")
	}
	offsets := make(map[Vector]struct{})
	for _, tile := range grid {
		if _, dupe := offsets[tile.Offset]; dupe {
			t.Fatal("GenTileGrid created duplicate offset")
		}
		offsets[tile.Offset] = struct{}{}
	}
	for x := 0; x < 20; x++ {
		for y := 0; y < 10; y++ {
			if _, ok := offsets[Vector{x, y}]; !ok {
				t.Fatal("GenTileGrid missing an offset")
			}
		}
	}
}

func TestGenTileGrid_Connectivity(t *testing.T) {
	const Cols, Rows = 20, 10
	grid := GenTileGrid(Cols, Rows, NewTile)
	for _, src := range grid {
		xBound := src.Offset.X == 0 || src.Offset.X == Cols-1
		yBound := src.Offset.Y == 0 || src.Offset.Y == Rows-1
		expectedLen := 8
		if xBound && yBound {
			expectedLen = 3
		} else if xBound || yBound {
			expectedLen = 5
		}
		if len(src.Adjacent) != expectedLen {
			t.Fatal("GenTileGrid gave incorrect number of links")
		}

		for delta, dst := range src.Adjacent {
			if dst.Adjacent[delta.Neg()] != src {
				t.Fatal("GenTileGrid failed to create backlink")
			}
		}
	}
}

func TestGenFence(t *testing.T) {
	const Cols, Rows = 20, 10
	grid := GenTileGrid(Cols, Rows, NewTile)
	GenFence(grid, func(t *Tile) {
		t.Pass = false
	})

	for _, tile := range grid {
		xBound := tile.Offset.X == 0 || tile.Offset.X == Cols-1
		yBound := tile.Offset.Y == 0 || tile.Offset.Y == Rows-1
		bound := xBound || yBound
		if bound && tile.Pass {
			t.Fatal("GenFence failed to modify boundary Tile")
		} else if !bound && !tile.Pass {
			t.Fatal("GenFence modified non-boundary Tile")
		}
	}
}
