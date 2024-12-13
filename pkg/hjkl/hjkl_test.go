package hjkl

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type MockGame struct {
	UpdateFn func([]Key) error
	DrawFn   func(Canvas)
}

func (m *MockGame) Update(ks []Key) error { return m.UpdateFn(ks) }
func (m *MockGame) Draw(c Canvas)         { m.DrawFn(c) }

type MockTerm struct {
	InitFn  func() error
	InputFn func() chan Key
	DoneFn  func()
	ClearFn func()
	BlitFn  func(Vector, Glyph)
	FlushFn func() error
}

func (m *MockTerm) Init() error            { return m.InitFn() }
func (m *MockTerm) Input() chan Key        { return m.InputFn() }
func (m *MockTerm) Done()                  { m.DoneFn() }
func (m *MockTerm) Clear()                 { m.ClearFn() }
func (m *MockTerm) Blit(v Vector, g Glyph) { m.BlitFn(v, g) }
func (m *MockTerm) Flush() error           { return m.FlushFn() }

func TestRun(t *testing.T) {
	// An incredibly fiddly test which could arguably benefit from a proper
	// mocking framework, but it does make sure that each tick has the correct
	// terminal calls in the correct order and that Run responds appropriately
	// to game updates.
	var log []string
	var update Key
	player := make(chan struct{})
	game := &MockGame{
		UpdateFn: func(ks []Key) error {
			if len(ks) == 0 {
				return nil
			}
			log = append(log, fmt.Sprintf("Update(%d)", ks[0]))
			for _, k := range ks {
				if k == KeyEsc {
					return Termination
				}
				update = k
			}
			return nil
		},
		DrawFn: func(c Canvas) {
			if update != 0 {
				log = append(log, "Draw")
				c.Blit(Vector{}, Glyph{})
			}
		},
	}
	term := &MockTerm{
		InitFn: func() error {
			log = append(log, "Init")
			return nil
		},
		InputFn: func() chan Key {
			log = append(log, "Input")
			c := make(chan Key)
			go func() {
				ks := []Key{'l', 'h', KeyEsc}
				for _, k := range ks {
					c <- k
					<-player
				}
			}()
			return c
		},
		DoneFn: func() { log = append(log, "Done") },
		ClearFn: func() {
			if update != 0 {
				log = append(log, "Clear")
			}
		},
		BlitFn: func(Vector, Glyph) {
			log = append(log, fmt.Sprintf("Blit(%d)", update))
		},
		FlushFn: func() error {
			if update != 0 {
				log = append(log, "Flush")
				update = 0
				player <- struct{}{}
			}
			return nil
		},
	}
	if err := Run(game, WithTerm(term), WithTPS(1000)); err != nil {
		t.Errorf("Run incorrectly gave error")
	}
	want := []string{
		"Init",
		"Input",
		"Update(108)",
		"Clear",
		"Draw",
		"Blit(108)",
		"Flush",
		"Update(104)",
		"Clear",
		"Draw",
		"Blit(104)",
		"Flush",
		"Update(27)",
		"Done",
	}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("Run call log is incorrect")
	}
}

func TestRun_InitError(t *testing.T) {
	want := errors.New("Init Error")
	game := &MockGame{}
	term := &MockTerm{
		InitFn: func() error { return want },
	}
	if got := Run(game, WithTerm(term)); got != want {
		t.Errorf("Run ignored init error")
	}
}

func TestRun_FlushError(t *testing.T) {
	want := errors.New("Flush Error")
	game := &MockGame{
		UpdateFn: func([]Key) error { return nil },
		DrawFn:   func(Canvas) {},
	}
	term := &MockTerm{
		InitFn:  func() error { return want },
		InputFn: func() chan Key { return nil },
		ClearFn: func() {},
		FlushFn: func() error { return want },
	}
	if err := Run(game, WithTerm(term)); err != want {
		t.Errorf("Run ignored flush error")
	}
}

func TestRun_UpdateError(t *testing.T) {
	want := errors.New("Update Error")
	game := &MockGame{
		UpdateFn: func([]Key) error { return want },
	}
	term := &MockTerm{
		InitFn:  func() error { return want },
		InputFn: func() chan Key { return nil },
		ClearFn: func() {},
	}
	if err := Run(game, WithTerm(term)); err != want {
		t.Errorf("Run ignored update error")
	}
}
