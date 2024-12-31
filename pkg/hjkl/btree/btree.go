// Package btree is an implementation of a behavior tree.
package btree

// State describes the outcome of running a Behavior.
type State int

// String converts a State to a string.
func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}

// State constants to be used by Behavior.
const (
	Unknown State = iota
	Running
	Success
	Failure
)

// Behavior is a node in a behavior tree.
type Behavior interface {
	Reset()
	Run() State
}

// Action is a function which acts as a Behavior.
type Action func() State

// Reset is a noop.
func (Action) Reset() {}

// Run calls the underlying function and returns the result.
func (f Action) Run() State {
	return f()
}

// Func is a function which acts as a Behavior which always succeeds.
type Func func()

// Reset is a noop.
func (Func) Reset() {}

// Run calls the underlying function and returns Success.
func (f Func) Run() State {
	f()
	return Success
}

// Conditional is a bool function which acts as a Behavior.
type Conditional func() bool

// Reset is a noop.
func (Conditional) Reset() {}

// Run calls the function and returns Success on true, or Failure otherwise.
func (f Conditional) Run() State {
	if f() {
		return Success
	}
	return Failure
}

// composite is the base of a Behavior composed of other Behavior.
type composite struct {
	nodes []Behavior
	index int
}

// Reset moves the index to 0 and resets all child Behavior.
func (c *composite) Reset() {
	c.index = 0
	for _, b := range c.nodes {
		b.Reset()
	}
}

// sequence is a Behavior which is he conjunction of child Behavior.
type sequence struct {
	composite
}

// Sequence gets a Behavior with the conjunction of child Behavior.
func Sequence(bs ...Behavior) Behavior {
	return &sequence{composite{nodes: bs}}
}

// Run runs each child Behavior in sequence. It succeeds if all the child
// Behavior succeed, but immediately fails if any child Behavior fails.
func (s *sequence) Run() State {
	for ; s.index < len(s.nodes); s.index++ {
		switch s.nodes[s.index].Run() {
		case Running:
			return Running
		case Success:
			continue
		case Failure:
			return Failure
		default:
			return Unknown
		}
	}
	return Success
}

// selection is a Behavior which is the disjunction of child Behavior.
type selection struct {
	composite
}

// Selection gets a Behavior with the disjunction of child Behavior.
func Selection(bs ...Behavior) Behavior {
	return &selection{composite{nodes: bs}}
}

// Run runs each child Behavior in sequence. It succeeds immediately if any the
// child Behavior succeed, but fails if all the child Behavior fails.
func (s *selection) Run() State {
	for ; s.index < len(s.nodes); s.index++ {
		switch s.nodes[s.index].Run() {
		case Running:
			return Running
		case Success:
			return Success
		case Failure:
			continue
		default:
			return Unknown
		}
	}
	return Failure
}

// decorator is a Behavior which transforms the output of another Behavior.
type decorator struct {
	node      Behavior
	transform func(State) State
}

// Reset resets the underlying Behavior.
func (d *decorator) Reset() {
	d.node.Reset()
}

// Run runs the underlying Behavior and returns the transformed result.
func (d *decorator) Run() State {
	return d.transform(d.node.Run())
}

// Invert wraps a Behavior to invert Success and Failure.
func Invert(b Behavior) Behavior {
	invert := func(s State) State {
		switch s {
		case Running:
			return Running
		case Success:
			return Failure
		case Failure:
			return Success
		default:
			return Unknown
		}
	}
	return &decorator{b, invert}
}

// Repeat wraps a Behavior to run indefinitely.
func Repeat(b Behavior) Behavior {
	repeat := func(s State) State {
		switch s {
		case Success, Failure:
			b.Reset()
			return Running
		case Running:
			return Running
		default:
			return Unknown
		}
	}
	return &decorator{b, repeat}
}

// ForceSuccess wraps a Behavior so Failure instead results in Success.
func ForceSuccess(b Behavior) Behavior {
	force := func(s State) State {
		switch s {
		case Success, Failure:
			return Success
		case Running:
			return Running
		default:
			return Unknown
		}
	}
	return &decorator{b, force}
}

// ForceFailure wraps a Behavior so Success instead results in Failure.
func ForceFailure(b Behavior) Behavior {
	force := func(s State) State {
		switch s {
		case Success, Failure:
			return Failure
		case Running:
			return Running
		default:
			return Unknown
		}
	}
	return &decorator{b, force}
}

// Until wraps a Behavior so it runs repeatedly until Success.
func Until(b Behavior) Behavior {
	until := func(s State) State {
		switch s {
		case Success:
			return Success
		case Failure:
			b.Reset()
			return Running
		case Running:
			return Running
		default:
			return Unknown
		}
	}
	return &decorator{b, until}
}

// While wraps a Behavior so it runs repeatedly until Failure.
func While(b Behavior) Behavior {
	while := func(s State) State {
		switch s {
		case Success:
			b.Reset()
			return Running
		case Failure:
			return Failure
		case Running:
			return Running
		default:
			return Unknown
		}
	}
	return &decorator{b, while}
}
