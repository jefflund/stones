package main

import "github.com/jefflund/stones/pkg/hjkl"

type Game struct{}

func (Game) Update(ks []hjkl.Key) error {
	if len(ks) > 0 {
		return hjkl.Termination
	}
	return nil
}

func (Game) Draw(c hjkl.Canvas) {
	for i, ch := range "Hello, World!" {
		c.Blit(hjkl.Vec(i, 0), hjkl.Ch(ch))
	}
}

func main() {
	if err := hjkl.Run(Game{}); err != nil {
		panic(err)
	}
}
