package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/shrtyk/pomodoro-cli/internal/app"
	cfg "github.com/shrtyk/pomodoro-cli/internal/config"
	p "github.com/shrtyk/pomodoro-cli/internal/player"
)

func main() {
	cfg := cfg.ParseConfig()

	player, err := p.NewPlayer(cfg.NotifyFile, cfg.DoneFile, cfg.NewRoundFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer player.Close()

	app, err := app.NewApplication(cfg, player)
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app.Start(ctx)
}
