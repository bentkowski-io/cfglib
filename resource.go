// package secretslib contains basic implementation of configuration reader.
package cfglib

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// NewWithMap returns a new instance of Configer from a map, like a map created from marshaling a json file.
func NewWithMapProvider(provider func() map[string]interface{}) *Config {
	return &Config{
		generic: provider(),
		opt:     json.JoinOptions(json.DefaultOptionsV2(), json.MatchCaseInsensitiveNames(true), json.StringifyNumbers(true)),
	}
}

// NewWithEnvProvider returns a new instance of Configer from environment variables and secrets provided by the provider func.
// The __ is used as a separator to define the hierarchy of the configuration.
func NewWithEnvProvider(envProvider func() []string) *Config {
	nodes := parseEnvVars(envProvider)
	nodes = flattenMap(nodes)
	cfgMap := &Config{
		generic: nodes,
		opt:     json.JoinOptions(json.DefaultOptionsV2(), json.MatchCaseInsensitiveNames(true), json.StringifyNumbers(true)),
	}

	return cfgMap
}

// New returns a new instance of Configer from environment variables and secrets.
// It takes environment function as provider to setup config with environemnt variables bound to runing container.
func New() *Config {
	return NewWithEnvProvider(os.Environ)
}

func NewFileJSONProvider(fileName string) (provider func() map[string]any, err error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bytes.NewBufferString("")
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	provider = func() map[string]any {
		cfg := make(map[string]any)
		err = json.UnmarshalDecode(jsontext.NewDecoder(buf), &cfg, json.DefaultOptionsV2())
		if err != nil {
			return nil
		}
		return flattenMap(cfg)
	}
	return
}

// NewFileEnvProvider returns a new provider that reads the environment variables from a file.
func NewFileEnvProvider(fileName string) (provider func() []string, err error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	lines := make([]string, 0)
	cs := ` '"`
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		if len(parts) != 2 {
			continue
		}
		parts[0] = strings.Trim(parts[0], cs)
		parts[1] = strings.Trim(parts[1], cs)
		if len(parts[0]) == 0 || strings.HasPrefix(parts[0], "#") {
			continue
		}

		lines = append(lines, strings.Join(parts, "="))
	}
	provider = func() []string {
		return lines
	}
	return
}
