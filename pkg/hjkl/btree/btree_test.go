package btree

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const Invalid State = 400

func TestAction(t *testing.T) {
	cases := []State{
		Unknown,
		Running,
		Success,
		Failure,
		Invalid,
	}
	for _, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			called := false
			got := Action(func() State {
				called = true
				return want
			}).Run()
			if got != want {
				t.Error("Action.Run gave incorrect state", got)
			}
			if !called {
				t.Error("Action.Run failed to call function")
			}
		})
	}
}

func TestFunc(t *testing.T) {
	called := false
	got := Func(func() {
		called = true
	}).Run()
	if got != Success {
		t.Error("Func.Run gave incorrect state", got)
	}
	if !called {
		t.Error("Func.Run failed to call function")
	}
}

func TestConditional(t *testing.T) {
	cases := []struct {
		Ret  bool
		Want State
	}{
		{true, Success},
		{false, Failure},
	}
	for _, c := range cases {
		t.Run(fmt.Sprint(c.Ret), func(t *testing.T) {
			called := false
			got := Conditional(func() bool {
				called = true
				return c.Ret
			}).Run()
			if got != c.Want {
				t.Error("Conditional.Run gave incorrect state", got)
			}
			if !called {
				t.Error("Conditional.Run failed to call function")
			}
		})
	}
}

func Recorded(states ...State) Behavior {
	i := 0
	return Action(func() State {
		if i >= len(states) {
			return Unknown
		}
		result := states[i]
		i++
		return result
	})
}

func BehaviorName(decorator any) string {
	full := runtime.FuncForPC(reflect.ValueOf(decorator).Pointer()).Name()
	parts := strings.Split(full, ".")
	return parts[len(parts)-1]
}

func ValidateComposite(t *testing.T, composite func(...Behavior) Behavior, children []Behavior, want []State) {
	name := BehaviorName(composite)
	b := composite(children...)
	got := make([]State, len(want))
	for i := range want {
		got[i] = b.Run()
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s produced incorrect states %v", name, got)
	}
}

func TestSequence_Success(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Success),
		Recorded(Running, Success),
		Recorded(Success),
	}
	want := []State{Running, Running, Success}
	ValidateComposite(t, Sequence, children, want)
}

func TestSequence_Failure(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Success),
		Recorded(Running, Running, Failure),
		Recorded(Success),
	}
	want := []State{Running, Running, Running, Failure}
	ValidateComposite(t, Sequence, children, want)
}

func TestSequence_Unknown(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Success),
		Recorded(Running, Running, Unknown),
		Recorded(Success),
	}
	want := []State{Running, Running, Running, Unknown}
	ValidateComposite(t, Sequence, children, want)
}

func TestSelection_Success(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
		Recorded(Success),
		Recorded(Success),
	}
	want := []State{Running, Running, Success}
	ValidateComposite(t, Selection, children, want)
}

func TestSelection_Failure(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
	}
	want := []State{Running, Running, Running, Failure}
	ValidateComposite(t, Selection, children, want)
}

func TestSelection_Unknown(t *testing.T) {
	children := []Behavior{
		Recorded(Running, Failure),
		Recorded(Failure),
		Recorded(Running, Failure),
		Recorded(Unknown),
		Recorded(Success),
	}
	want := []State{Running, Running, Unknown}
	ValidateComposite(t, Selection, children, want)
}

type MockBehavior struct {
	RunFn   func() State
	ResetFn func()
}

func (m *MockBehavior) Run() State { return m.RunFn() }

func (m *MockBehavior) Reset() { m.ResetFn() }

func VerifyDecorator(t *testing.T, decorator func(Behavior) Behavior, mapping map[State]State) {
	name := BehaviorName(decorator)
	for ret, want := range mapping {
		t.Run(fmt.Sprint(ret), func(t *testing.T) {
			called := false
			reset := false
			a := decorator(&MockBehavior{
				RunFn: func() State {
					called = true
					return ret
				},
				ResetFn: func() {
					reset = true
				},
			})
			if got := a.Run(); got != want {
				t.Errorf("%s.Run gave incorrect state %v", name, got)
			}
			if !called {
				t.Errorf("%s.Run failed to run Behavior", name)
			}
			a.Reset()
			if !reset {
				t.Errorf("%s.Reset  failed to reset Behavior", name)
			}
		})
	}
}

func TestInvert(t *testing.T) {
	VerifyDecorator(t, Invert, map[State]State{
		Running: Running,
		Success: Failure,
		Failure: Success,
		Unknown: Unknown,
		Invalid: Unknown,
	})
}

func TestForceSuccess(t *testing.T) {
	VerifyDecorator(t, ForceSuccess, map[State]State{
		Running: Running,
		Success: Success,
		Failure: Success,
		Unknown: Unknown,
		Invalid: Unknown,
	})
}

func TestForceFailure(t *testing.T) {
	VerifyDecorator(t, ForceFailure, map[State]State{
		Running: Running,
		Success: Failure,
		Failure: Failure,
		Unknown: Unknown,
		Invalid: Unknown,
	})
}

type RepeatedDecoratorCase struct {
	Ret       State
	WantState State
	WantReset bool
}

func VerifyRepeatedDecorator(t *testing.T, decorator func(Behavior) Behavior, cases []RepeatedDecoratorCase) {
	name := BehaviorName(decorator)
	var ret State
	var reset bool
	a := decorator(&MockBehavior{
		RunFn: func() State {
			return ret
		},
		ResetFn: func() {
			reset = true
		},
	})
	for i, c := range cases {
		ret = c.Ret
		reset = false
		if got := a.Run(); got != c.WantState {
			t.Errorf("%s.Run gave incorrect state %v on step %d", name, got, i)
		}
		if reset != c.WantReset {
			t.Errorf("%s.Run gave incorrect reset on step %d", name, i)
		}
	}
}

func TestRepeat(t *testing.T) {
	VerifyRepeatedDecorator(t, Repeat, []RepeatedDecoratorCase{
		{Success, Running, true},
		{Success, Running, true},
		{Failure, Running, true},
		{Failure, Running, true},
		{Running, Running, false},
		{Running, Running, false},
		{Unknown, Unknown, false},
		{Invalid, Unknown, false},
	})
}

func TestUntil(t *testing.T) {
	VerifyRepeatedDecorator(t, Until, []RepeatedDecoratorCase{
		{Failure, Running, true},
		{Running, Running, false},
		{Failure, Running, true},
		{Running, Running, false},
		{Success, Success, false},
		{Unknown, Unknown, false},
		{Invalid, Unknown, false},
	})
}

func TestWhile(t *testing.T) {
	VerifyRepeatedDecorator(t, While, []RepeatedDecoratorCase{
		{Success, Running, true},
		{Running, Running, false},
		{Success, Running, true},
		{Running, Running, false},
		{Failure, Failure, false},
		{Unknown, Unknown, false},
		{Invalid, Unknown, false},
	})
}
