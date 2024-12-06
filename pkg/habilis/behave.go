package habilis

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

type ActTrigger struct {
	Keys []hjkl.Key
}

func Contoller(m *rl.Mob, v *ActTrigger) {
	for _, k := range v.Keys {
		if delta, ok := hjkl.VIKeyMap[k]; ok {
			m.Move(delta)
		}
	}
}

func Wander(m *rl.Mob, v *ActTrigger) {
	m.Move(rand.Choice(hjkl.Dirs8))
}
