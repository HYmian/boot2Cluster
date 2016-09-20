package conf

import (
	"testing"
)

func TestAddNode(t *testing.T) {
	cfg, err := LoadConfig("conf.yml")
	if err != nil {
		t.Errorf("load config error: %s", err.Error())
	}

	boot := NewBoot(cfg, 3)
	if err := boot.AddNode("m1", "192", 1); err != nil {
		t.Error(err.Error())
	}
	if err := boot.AddNode("m2", "193", 2); err != nil {
		t.Error(err.Error())
	}
	if err := boot.AddNode("m3", "194", 3); err != nil {
		t.Error(err.Error())
	}
}
