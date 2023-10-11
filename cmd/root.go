package cmd

import (
	"context"
	"log/slog"

	"github.com/bep/simplecobra"
)

type RootCommand struct {
	log    *slog.Logger
	logLvl *slog.LevelVar
}

func NewRootCommand(logLevelVar *slog.LevelVar) *RootCommand {
	return &RootCommand{
		log:    slog.Default(),
		logLvl: logLevelVar,
	}
}

func (r *RootCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{
		newExecCommand(),
		newInitCommand(),
	}
}

func (r *RootCommand) Init(c *simplecobra.Commandeer) error {
	c.CobraCommand.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging for debugging purposes")
	return nil
}

func (r *RootCommand) Name() string {
	return "guardian-cli"
}

func (r *RootCommand) PreRun(c *simplecobra.Commandeer, c1 *simplecobra.Commandeer) error {
	level := slog.LevelWarn

	if vFlag, _ := c.CobraCommand.Flags().GetBool("verbose"); vFlag == true {
		level = slog.LevelDebug
	}
	r.logLvl.Set(level)
	r.log = slog.Default().With("command", "root")
	return nil
}

func (r *RootCommand) Run(ctx context.Context, c *simplecobra.Commandeer, s []string) error {
	return nil
}
