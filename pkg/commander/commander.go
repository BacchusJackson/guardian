// Package commander is used
//
// Example Usage
package commander

import (
	"encoding/json"
	"errors"
	"io"
	"text/template"

	"log/slog"
)

// Command is a data representation of a single command
type Command struct {
	Description string            `json:"description,omitempty"`
	Template    string            `json:"template"`
	Values      map[string]string `json:"values"`
}

// Print parses the Template field and executes uses the Values field
func (c *Command) Print(w io.Writer) error {
	tmpl, err := template.New("command").Parse(c.Template)
	if err != nil {
		slog.Error("failed to parse template", "err", err)
		return err
	}
	return tmpl.Execute(w, c.Values)
}

// CommandMap is a map of multiple commands
type CommandMap map[string]*Command

// ParsePrint creates a command and prints to the io.Writer
func ParsePrint(w io.Writer, template string, values map[string]string) error {
	cmd := &Command{
		Template: template,
		Values:   values,
	}
	return cmd.Print(w)
}

// ParsePrintFrom decodes the sorce config and parses only the target command by name
func ParsePrintFrom(srcConfig io.Reader, dest io.Writer, targetCmd string) error {
	config := make(CommandMap, 0)
	err := json.NewDecoder(srcConfig).Decode(&config)
	if err != nil {
		slog.Error("failed decoding", "err", err)
		return err
	}
	cmd, ok := config[targetCmd]
	if !ok {
		slog.Error("target command not found", "target", targetCmd)
		return errors.New("input")
	}
	return cmd.Print(dest)
}

func PrintExampleJSON(w io.Writer) error {
	cmd := CommandMap{
		"echo": &Command{
			Description: "print a text message",
			Template:    "echo {{- if .newline}} -n {{end -}} \"{{.msg}}\"",
			Values:      map[string]string{"msg": "howdy world", "newline": "true"},
		},
		"docker-build": &Command{
			Description: "build an image with docker",
			Template:    "docker {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}",
			Values:      map[string]string{"file": "Dockerfile.custom"},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(&cmd)
}
