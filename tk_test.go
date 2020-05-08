package tk_test

import (
	"os"
	"testing"

	"github.com/assay-it/tk"
)

func TestEnvDefault(t *testing.T) {
	if tk.Env("VAR1", "undef") != "undef" {
		t.Error("VAR1 is defined")
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("VAR2", "defined")

	if tk.Env("VAR2", "undef") != "defined" {
		t.Error("VAR is not defined")
	}
}
