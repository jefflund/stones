package clock

import (
	"reflect"
	"sort"
	"testing"
)

func TestClock(t *testing.T) {
	actions := map[string][]int{
		"a": {1, 2, 3, 3, 2, 1},
		"b": {3, 3, 3, 3, 3},
		"c": {5, 2},
	}
	want := [][]string{
		{"a"},
		nil,
		{"a", "b"},
		nil,
		{"c"},
		{"a", "b"},
		{"c"},
		nil,
		{"a", "b"},
		nil,
		{"a"},
		{"a", "b"},
	}

	c := New[string]()
	schedule := func(tick []string) {
		for _, s := range tick {
			if as := actions[s]; len(as) > 0 {
				c.Schedule(s, as[0])
				actions[s] = as[1:]
			}
		}
	}
	for s := range actions {
		schedule([]string{s})
	}

	var got [][]string
	for range want {
		tick := c.Tick()
		schedule(tick)
		sort.Strings(tick)
		got = append(got, tick)
	}

	if !reflect.DeepEqual(want, got) {
		t.Error("Clock ticked incorrectly. Got:", got)
	}
}

func TestClock_Unschedule(t *testing.T) {
	c := New[string]()
	c.Schedule("a", 1)
	c.Schedule("a", 2)
	c.Schedule("b", 2)
	c.Schedule("c", 1)
	c.Unschedule("b")
	if !reflect.DeepEqual(c.Tick(), []string{"c"}) {
		t.Error("Clock.Schedule failed to reschedule")
	}
	if !reflect.DeepEqual(c.Tick(), []string{"a"}) {
		t.Error("Clock.Unschedule failed to unschedule")
	}
}
