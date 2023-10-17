// Package commander is used
//
// Example Usage
package commander

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
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
	var err error
	tmpl := template.New("command")
	tmpl = tmpl.Funcs(template.FuncMap{"mflag": multiFlag})
	tmpl, err = tmpl.Parse(c.Template)
	if err != nil {
		slog.Error("failed to parse template", "err", err)
		return err
	}
	return tmpl.Execute(w, c.Values)
}

// indirect returns the item at the end of indirection, and a bool to indicate
// if it's nil. If the returned bool is true, the returned value's kind will be
// either a pointer or interface.
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}

// multiFlag ...
func multiFlag(item reflect.Value) (string, error) {
	item, isNil := indirect(item)
	if isNil {
		return "", fmt.Errorf("args of nil pointer")
	}
	if item.Kind() != reflect.String {
		return "", fmt.Errorf("args of type %s", item.Type())
	}

	inputString := item.String()
	splitChar := []rune(inputString)[0]

	parts := strings.Split(inputString, string(splitChar))
	if len(parts) < 3 {
		return "", fmt.Errorf("got: %s want format like: --arg \"value 1\" \"value 2\"", inputString)
	}

	flagString := parts[1]
	out := make([]string, 0, len(parts)*2)

	// 0 is Blank since the first character is the split character
	// 1 is the flag arg
	for _, part := range parts[2:] {
		out = append(out, flagString, part)
	}

	return strings.Join(out, " "), nil
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
		Template:    "docker build {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}",
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
