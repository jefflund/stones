package main

import "github.com/jefflund/stones/pkg/hjkl"

func main() {
	term := hjkl.TermboxTerminal{}

	if err := term.Init(); err != nil {
		panic(err)
	}
	defer term.Done()

	term.Clear()
	for i, ch := range "Hello, World!" {
		term.Blit(hjkl.Vec(i, 0), hjkl.Ch(ch))
	}
	term.Flush()

	<-term.Input()
}
