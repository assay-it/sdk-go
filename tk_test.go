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

func TestHostBuildDefault(t *testing.T) {
	for _, h := range []string{
		"localhost",
		"localhost:8080",
	} {
		if host := tk.Host(h); host != h {
			t.Errorf("Unexpected default host: %s", host)
		}
	}
}

func TestHostBuildEnv(t *testing.T) {
	for in, ex := range map[string]string{
		"https://localhost":                  "localhost",
		"https://localhost:8080":             "localhost:8080",
		"https://localhost.localdomain":      "localhost.localdomain",
		"https://localhost.localdomain:8080": "localhost.localdomain:8080",
	} {
		os.Setenv("BUILD_ENDPOINT", in)
		if host := tk.Host(""); host != ex {
			t.Errorf("Unexpected host environment %s", host)
		}
	}
	os.Unsetenv("BUILD_ENDPOINT")
}

func TestHostBuildComposition(t *testing.T) {
	os.Setenv("BUILD_ID", "0")
	os.Setenv("BUILD_DOMAIN", "example.com")

	if host := tk.Host(""); host != "v0.example.com" {
		t.Errorf("Unexpected host composition %s", host)
	}
	os.Unsetenv("BUILD_ID")
	os.Unsetenv("BUILD_DOMAIN")
}
