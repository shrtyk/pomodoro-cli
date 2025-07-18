package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shrtyk/pomodoro-cli/internal/config"
	"github.com/shrtyk/pomodoro-cli/internal/player"
)

// Application defines the interface for the Pomodoro application.
type Application interface {
	Start(ctx context.Context)
}

// Phase represents the current state of the Pomodoro timer.
type Phase int

const (
	_ Phase = iota
	round
	notify
	rest
)

// appState holds the dynamic state of a running Pomodoro timer.
type appState struct {
	app       *application
	wg        *sync.WaitGroup
	phase     Phase
	timer     *time.Timer
	switchers map[Phase]func() bool
}

// application holds the configuration and state of the Pomodoro timer.
type application struct {
	rounds         uint64
	roundDuration  time.Duration
	restDuration   time.Duration
	notifyDuration time.Duration
	player         player.Player
	state          *appState
}

// NewApplication creates a new Application instance.
func NewApplication(cfg *config.Config, player player.Player) (Application, error) {
	app := &application{
		rounds:         cfg.Rounds,
		roundDuration:  cfg.RoundDuration,
		restDuration:   cfg.RestDuration,
		notifyDuration: cfg.NotifyDuration,
		player:         player,
	}

	state := &appState{
		app:   app,
		wg:    &sync.WaitGroup{},
		phase: round,
	}

	state.switchers = map[Phase]func() bool{
		round:  state.roundState,
		notify: state.notifyState,
		rest:   state.restState,
	}

	app.state = state
	return app, nil
}

// switchState handles the transition between different states.
// It returns true if all rounds are completed.
func (app *application) switchState() (done bool) {
	app.state.wg.Add(1)
	op, ok := app.state.switchers[app.state.phase]
	if !ok {
		// This should not happen in normal operation.
		panic("no such state")
	}
	return op()
}

// roundState is the handler for the 'round' state.
func (s *appState) roundState() (done bool) {
	go func() {
		defer s.wg.Done()
		s.app.player.PlayLoopAlert()
	}()
	s.app.rounds--
	if s.app.rounds == 0 {
		done = true
		return
	}
	s.phase = notify
	s.timer.Reset(time.Until(time.Now().Add(s.app.restDuration - s.app.notifyDuration)))
	return
}

// notifyState is the handler for the 'notify' state.
func (s *appState) notifyState() (done bool) {
	go func() {
		defer s.wg.Done()
		s.app.player.PlaySoonAlert()
	}()
	s.phase = rest
	s.timer.Reset(s.app.notifyDuration)
	return
}

// restState is the handler for the 'rest' state.
func (s *appState) restState() (done bool) {
	go func() {
		defer s.wg.Done()
		s.app.player.PlayRoundAlert()
	}()
	s.phase = round
	s.timer.Reset(s.app.roundDuration)
	return
}

// Start begins the Pomodoro timer.
// It initializes the timer for precision and listens for context cancellation
// for graceful shutdown.
func (app *application) Start(ctx context.Context) {
	app.state.timer = time.NewTimer(app.roundDuration)
	app.state.wg.Add(1)

	go func() {
		defer app.state.wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Got a syscall. Shutting down...")
				// Stop the timer to prevent the goroutine from leaking.
				if !app.state.timer.Stop() {
					select {
					case <-app.state.timer.C:
					default:
					}
				}
				return
			case <-app.state.timer.C:
				if app.switchState() {
					fmt.Println("DONE")
					return
				}
			}
		}
	}()

	app.state.wg.Wait()
}
