package cfglib

import (
	"fmt"
	"strings"

	json "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Converter[T any] struct {
	m         *Config
	converter func(a any) T
}

type Config struct {
	generic map[string]any
	opt     json.Options
}

func (m *Config) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range m.generic {
		keys = append(keys, k)
	}
	return keys
}

// Read returns raw value of key. It can either be string or a map[string]any.
func (m *Config) Read(key string) any {
	return m.generic[(ToSnakeCase(key))]
}

// Read returns raw value of key. It can either be string or a map[string]any.
func (m *Config) ReadString(key string) string {
	if v, ok := m.generic[(ToSnakeCase(key))]; ok {
		return fmt.Sprint(v)
	} else {
		return ""
	}
}

func (m *Config) ReadSlice(key string, cfg interface{}) error {
	v := m.Read(ToSnakeCase(key))
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.UnmarshalDecode(jsontext.NewDecoder(strings.NewReader(string(b))), &cfg, m.opt)
}

// ReadMap parses a value into a map[string]any pointer.
func (m *Config) ReadMap(key string, cfg any) error {
	v := m.Read(ToSnakeCase(key))
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = json.UnmarshalDecode(jsontext.NewDecoder(strings.NewReader(string(b))), &cfg, m.opt)
	if err != nil {
		return err
	}
	return nil
}

// ReadMap parses a value into a struct pointer.
func (m *Config) ReadStruct(key string, cfg any) error {
	v := m.Read(ToSnakeCase(key))
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	json.UnmarshalDecode(jsontext.NewDecoder(strings.NewReader(string(b))), cfg, m.opt)
	return json.Unmarshal(b, &cfg)
}

func (c Converter[T]) Read(key string) T {
	return c.converter(c.m.Read(key))
}

// ToSnakeCase consverts string to snake case
func ToSnakeCase(s string) (result string) {
	var previous rune
	chageDone := false
	for i, v := range s {
		if !chageDone && i > 0 && caseChanged(v, previous) {
			result += "_"
		}

		chageDone = false
		result += strings.ToLower(string(v))
		previous = v
	}
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, Separator, ".")

	return strings.ToLower(result)

}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func caseChanged(currrent, previous rune) bool {
	if isSpecialChar(currrent) || isSpecialChar(previous) {
		return false
	}
	return isUpper(currrent) && !isUpper(previous)
}

func isSpecialChar(r rune) bool {
	if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
		return false
	}
	return true
}
