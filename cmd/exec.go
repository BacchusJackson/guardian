package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/bep/simplecobra"
	"github.com/gatecheckdev/guardian/pkg/commander"
)

type execCommand struct {
	log        *slog.Logger
	dryRun     bool
	values     map[string]string
	template   string
	configFile *os.File
	target     string
}

func newExecCommand() *execCommand {
	return &execCommand{log: slog.Default().With("command", "exec")}
}
func (_ *execCommand) Commands() []simplecobra.Commander {
	return make([]simplecobra.Commander, 0)
}

func (_ *execCommand) Init(c *simplecobra.Commandeer) error {
	c.CobraCommand.Short = "Execute a command with options or from a configuration file"
	c.CobraCommand.Flags().BoolP("dry-run", "n", false, "Print the command without executing")
	c.CobraCommand.Flags().StringToStringP("values", "i", make(map[string]string), "Key value pairs for the template")
	c.CobraCommand.Flags().StringP("template", "t", "", "The template string for the command")

	c.CobraCommand.Flags().StringP("file", "f", "commander.json", "The configuration file with key-values and templates")
	c.CobraCommand.Flags().String("target", "", "The name of the target command")

	c.CobraCommand.MarkFlagFilename("file")
	c.CobraCommand.MarkFlagsRequiredTogether("file", "target")
	return nil
}

func (_ *execCommand) Name() string {
	return "exec"
}

func (e *execCommand) PreRun(c *simplecobra.Commandeer, c1 *simplecobra.Commandeer) error {
	dryRunFlag, _ := c.CobraCommand.Flags().GetBool("dry-run")
	e.dryRun = dryRunFlag
	valuesFlag, _ := c.CobraCommand.Flags().GetStringToString("values")
	e.values = valuesFlag
	templateFlag, _ := c.CobraCommand.Flags().GetString("template")
	e.template = templateFlag
	targetFlag, _ := c.CobraCommand.Flags().GetString("target")
	e.target = targetFlag

	fileFlag, _ := c.CobraCommand.Flags().GetString("file")
	if fileFlag == "" {
		return nil
	}
	slog.Debug("open config file", "filename", fileFlag)
	f, err := os.Open(fileFlag)
	if err != nil {
		slog.Error("failed to open config file", "filename", fileFlag)
		return err
	}

	e.configFile = f
	return nil
}

func (e *execCommand) Run(ctx context.Context, c *simplecobra.Commandeer, s []string) error {
	e.log.Debug("", "values", e.values, "template", e.template, "target", e.target, "config_file", e.configFile)
	commander.ParsePrint(c.CobraCommand.OutOrStdout(), e.template, e.values)

	switch {
	case e.dryRun && e.configFile != nil:
		return commander.ParsePrintFrom(e.configFile, c.CobraCommand.OutOrStdout(), e.target)
	case e.dryRun && e.template != "":
		return commander.ParsePrint(c.CobraCommand.OutOrStdout(), e.template, e.values)
	case !e.dryRun:
		e.log.Warn("executing commands directly is not yet supported. Please use `--dry-run`")
		return errors.New("not supported")
	default:
		e.log.Warn("unexpected input")
		return errors.New("no config file or template argument")
	}
}
