package loglevel

import (
	"context"
	"os"

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

func FromEnv(prefix string) zerolog.Level {
	return FromLookuper(prefix, envconfig.OsLookuper())
}

func FromLookuper(prefix string, l envconfig.Lookuper) zerolog.Level {
	l = envconfig.PrefixLookuper(prefix, l)
	c := getC(l)

	if c.Debug {
		return zerolog.DebugLevel
	}

	return zerolog.InfoLevel
}

func getC(l envconfig.Lookuper) Config {
	var c Config
	if err := envconfig.ProcessWith(context.Background(), &c, l); err != nil {
		// do nothing
	}
	return c
}

func Baseline(prefix string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stderr).With().
		Timestamp().
		Logger().
		Level(FromEnv(prefix))
}
