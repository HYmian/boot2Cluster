package conf

import (
	"testing"
)

func TestExec(t *testing.T) {
	err := Exec("echo hehe")
	if err != nil {
		t.Errorf("test exec error: %s", err.Error())
	}
}
