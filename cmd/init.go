package cmd

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/bep/simplecobra"
	"github.com/gatecheckdev/guardian/pkg/commander"
)

type initCommand struct {
	log            *slog.Logger
	outputFiletype string
}

func newInitCommand() *initCommand {
	return &initCommand{log: slog.Default().With("command", "init")}
}
func (_ *initCommand) Commands() []simplecobra.Commander {
	return make([]simplecobra.Commander, 0)
}

func (_ *initCommand) Init(c *simplecobra.Commandeer) error {
	c.CobraCommand.Short = "Print an example commander file to standard out"

	c.CobraCommand.Flags().StringP("output", "o", "json", "output filetype ['json','yaml','toml']")
	return nil
}

func (_ *initCommand) Name() string {
	return "init"
}

func (e *initCommand) PreRun(c *simplecobra.Commandeer, c1 *simplecobra.Commandeer) error {
	o, _ := c.CobraCommand.Flags().GetString("output")
	e.outputFiletype = o
	return nil
}

func (e *initCommand) Run(ctx context.Context, c *simplecobra.Commandeer, s []string) error {
	slog.Debug("run", "output_filetype", e.outputFiletype)
	switch strings.TrimSpace(strings.ToLower(e.outputFiletype)) {
	case "json":
		return commander.PrintExampleJSON(c.CobraCommand.OutOrStdout())
	case "yaml":
		return commander.PrintExampleYAML(c.CobraCommand.OutOrStdout())
	case "toml":
		return commander.PrintExampleTOML(c.CobraCommand.OutOrStdout())
	}
	return errors.New("unsupported filetype")
}
