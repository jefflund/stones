package gen

import (
	"testing"

	"github.com/jefflund/stones/pkg/hjkl"
)

func TestGenTileGrid(t *testing.T) {
	const W, H = 10, 5
	tiles := GenTileGrid(W, H, hjkl.NewTile)

	if len(tiles) != W*H {
		t.Error("GenTileGrid return incorrect number of Tile")
	}

	for _, src := range tiles {
		xBound := src.Offset.X == 0 || src.Offset.X == W-1
		yBound := src.Offset.Y == 0 || src.Offset.Y == H-1
		want := 8
		if xBound && yBound {
			want = 3
		} else if xBound || yBound {
			want = 5
		}
		if len(src.Adjacent) != want {
			t.Error("GenTileGrid got incorrect number of links")
		}

		for delta, dst := range src.Adjacent {
			if dst.Adjacent[delta.Neg()] != src {
				t.Error("GenTileGrid failed to create backlink")
			}
		}
	}
}

func TestGenFence(t *testing.T) {
	const W, H = 10, 5

	tiles := GenTileGrid(W, H, hjkl.NewTile)
	GenFence(tiles, func(t *hjkl.Tile) {
		t.Pass = false
	})

	for _, src := range tiles {
		bound := src.Offset.X == 0 || src.Offset.X == W-1 || src.Offset.Y == 0 || src.Offset.Y == H-1
		if bound && src.Pass {
			t.Error("GenFence failed to modify boundary Tile")
		} else if !bound && !src.Pass {
			t.Error("GenFence modified non-boundary Tile")
		}
	}
}
