package loglevel

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Debug bool `env:"DEBUG,default=false"`
}

func SetGlobal() {
	var c Config
	if err := envconfig.Process(context.Background(), &c); err != nil {
		// do nothing, leave it at the default "not debug, therefore info+"
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if c.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func SetLoggerLevel(l zerolog.Logger) zerolog.Logger {
	lvl := zerolog.InfoLevel
	var c Config
	if err := envconfig.Process(context.Background(), &c); err != nil {
		// do nothing, leave it at the default "not debug, therefore info+"
	}

	if c.Debug {
		lvl = zerolog.DebugLevel
	}

	return l.Level(lvl)
}
