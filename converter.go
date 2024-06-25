package cfglib

import (
	"fmt"
	"strconv"
)

func NewBoolConverter(m *Config) Converter[bool] {
	return Converter[bool]{
		m: m,
		converter: func(a any) bool {
			if a == nil {
				return false
			}
			b, _ := strconv.ParseBool(fmt.Sprint(a))
			return b
		},
	}
}
func NewIntConverter(m *Config) Converter[int] {
	return Converter[int]{
		m: m,
		converter: func(a any) int {
			if a == nil {
				return 0
			}
			v, _ := strconv.Atoi(fmt.Sprint(a))
			return v
		},
	}
}
