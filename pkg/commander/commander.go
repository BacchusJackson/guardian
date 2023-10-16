// Package commander is used
//
// Example Usage
package commander

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"text/template"

	"log/slog"

	"github.com/sclevine/yj/convert"
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

type FileExt string

const (
	ExtJSON  FileExt = "json"
	ExtYAML          = "yaml"
	ExtTOML          = "toml"
	ExtOther         = ""
)

// ParsePrintFrom decodes the sorce config and parses only the target command by name
func ParsePrintFrom(extType FileExt, srcConfig io.Reader, dest io.Writer, targetCmd string) error {

	buf := new(bytes.Buffer)
	config := make(CommandMap, 0)
	var readerErr error

	switch extType {
	case ExtJSON:
		slog.Debug("read json")
		_, err := buf.ReadFrom(srcConfig)
		readerErr = err
	case ExtYAML:
		slog.Debug("decode yaml")
		obj, err := convert.YAML{}.Decode(srcConfig)
		readerErr = err
		convert.JSON{}.Encode(buf, obj)

	case ExtTOML:
		slog.Debug("decode toml")
		obj, err := convert.TOML{}.Decode(srcConfig)
		readerErr = err
		convert.JSON{}.Encode(buf, obj)
	case ExtOther:
		return errors.New("file extension is not JSON/YAML/TOML")
	}
	if readerErr != nil {
		slog.Error("there's a problem with the input reader", "err", readerErr, "ext_type", extType)
		return readerErr
	}

	if err := json.NewDecoder(buf).Decode(&config); err != nil {
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

var exampleCmd = CommandMap{
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

func JSONtoYAML(jsonIn io.Reader, yamlOut io.Writer) error {
	obj, err := convert.JSON{}.Decode(jsonIn)
	if err != nil {
		slog.Error("failed to decode JSON", "err", err)
		return err
	}

	return convert.YAML{}.Encode(yamlOut, obj)
}

func JSONtoTOML(jsonIn io.Reader, tomlOut io.Writer) error {
	obj, err := convert.JSON{}.Decode(jsonIn)
	if err != nil {
		slog.Error("failed to decode JSON", "err", err)
		return err
	}
	return new(convert.TOML).Encode(tomlOut, obj)
}

func PrintExampleJSON(w io.Writer) error {
	return convert.JSON{EscapeHTML: false, Indent: true}.Encode(w, &exampleCmd)
}

func PrintExampleYAML(w io.Writer) error {

	pr, pw := io.Pipe()
	defer pr.Close()

	go func() {
		defer pw.Close()
		err := convert.JSON{}.Encode(pw, &exampleCmd)
		if err != nil {
			slog.Error("failed to encode JSON", "err", err)
		}
	}()

	return JSONtoYAML(pr, w)
}

func PrintExampleTOML(w io.Writer) error {

	pr, pw := io.Pipe()
	defer pr.Close()

	go func() {
		defer pw.Close()
		err := convert.JSON{}.Encode(pw, &exampleCmd)
		if err != nil {
			slog.Error("failed to encode JSON", "err", err)
		}
	}()

	return JSONtoTOML(pr, w)
}
