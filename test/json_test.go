package cfglib_test

import (
	"testing"

	"github.com/bentkowski-io/cfglib"
)

func prepareJSON(t *testing.T) func() map[string]any {
	provider, err := cfglib.NewFileJSONProvider("cfg.json")
	if err != nil {
		t.Fatal(err)
	}
	return provider
}

func TestJSON(t *testing.T) {

	cfg := cfglib.NewWithMapProvider(prepareJSON(t))
	services := Services{}
	if err := cfg.ReadMap("srv", &services); err != nil {
		t.Fatal(err)
	}

}
