package conf

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("conf.yml")
	if err != nil {
		t.Errorf("load config error: %s", err.Error())
	}

	t.Log(cfg)
}
