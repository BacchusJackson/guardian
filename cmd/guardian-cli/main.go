package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gatecheckdev/guardian/cmd"
	"github.com/lmittmann/tint"
)

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelWarn)

	h := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      lvl,
		TimeFormat: time.TimeOnly,
	})
	slog.SetDefault(slog.New(h))

	cmd, err := simplecobra.New(cmd.NewRootCommand(lvl))
	if err != nil {
		slog.Error("failed to generate command", "err", err)
		os.Exit(1)
	}

	if _, err = cmd.Execute(context.Background(), os.Args[1:]); err != nil {
		slog.Error("command execution failure", "err", err)
		os.Exit(2)
	}
}
