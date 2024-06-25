package cfglib_test

import (
	"strings"
	"testing"

	"bentkowski.io/cfglib"
)

func TestCompare(t *testing.T) {
	cfg1 := cfglib.NewWithEnvProvider(prepareEnv(t))
	cfg2 := cfglib.NewWithMapProvider(prepareJSON(t))

	for _, key := range cfg1.Keys() {
		t.Run(key, func(t *testing.T) {
			if !strings.EqualFold(cfg1.ReadString(key), cfg2.ReadString(key)) {
				t.Errorf("expected %s, got %s", cfg2.ReadString(key), cfg1.ReadString(key))
			}
		})
	}
}
