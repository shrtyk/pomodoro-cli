package config

import (
	"flag"
	"time"
)

type Config struct {
	Rounds         uint64
	RoundDuration  time.Duration
	RestDuration   time.Duration
	NotifyDuration time.Duration

	NotifyFile   string
	DoneFile     string
	NewRoundFile string
}

func ParseConfig() *Config {
	cfg := new(Config)

	flag.Uint64Var(&cfg.Rounds, "rounds", 1, "Rounds amount")
	flag.DurationVar(&cfg.RoundDuration, "round_duration", 25*time.Minute, "Duration of one round")
	flag.DurationVar(&cfg.RestDuration, "rest_duration", 5*time.Minute, "Duration of the rest")
	flag.DurationVar(&cfg.NotifyDuration, "notify_before_rest_end", 30*time.Second, "Notify before end of the rest")

	flag.StringVar(&cfg.NotifyFile, "notification_file", "notify.mp3", "Path to the file to play as notification before end of the rest")
	flag.StringVar(&cfg.DoneFile, "done_file", "done.mp3", "Path to the file to play at the end of the rest")
	flag.StringVar(&cfg.NewRoundFile, "round_file", "round.mp3", "Path to the file to play at the start of the round")

	flag.Parse()
	return cfg
}
