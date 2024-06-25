package cfglib_test

import (
	"testing"

	"bentkowski.io/cfglib"
)

// Services contains the url to access the web services.
type Services map[string]Host

// Host contains the url and the id of the application.
// string option in json tag is used to destring a string id into int64 during json unmarshalling.
type Host struct {
	Host     string `json:"host"`
	TargetID int64  `json:"id,string"`
}

// TestServices tests unmarshaling services struct from json.
// The services contains Host struct, where the id is converted to string in json but int64 in struct.
// The test should pass, due to using string option in json tag.
func TestServices(t *testing.T) {
	cfg := cfglib.NewWithEnvProvider(prepareEnv(t))
	services := Services{}
	if err := cfg.ReadMap("srv", &services); err != nil {
		t.Fatal(err)
	}

}

func TestSlice(t *testing.T) {
	cfg := cfglib.NewWithEnvProvider(prepareEnv(t))
	var slice []int
	if err := cfg.ReadSlice("slice", &slice); err != nil {
		t.Fatal(err)
	}
}

func TestTrim(t *testing.T) {
	cfg := cfglib.NewWithEnvProvider(prepareEnv(t))
	test := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d=",
	}
	for k, v := range test {
		t.Run(k, func(t *testing.T) {
			if cfg.ReadString(k) != v {
				t.Errorf("expected %s, got %s", v, cfg.ReadString(k))
			}
		})
	}
}

func prepareEnv(t *testing.T) func() []string {
	lFunc, err := cfglib.NewFileEnvProvider(".env")
	if err != nil {
		t.Fatal(err)
	}
	return lFunc
}
