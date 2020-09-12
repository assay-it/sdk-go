package assay_test

import (
	"os"
	"testing"

	"github.com/assay-it/sdk-go/assay"
)

func TestEnvDefault(t *testing.T) {
	if assay.Env("VAR1", "undef") != "undef" {
		t.Error("VAR1 is defined")
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("VAR2", "defined")

	if assay.Env("VAR2", "undef") != "defined" {
		t.Error("VAR is not defined")
	}
}

func TestHostBuildDefault(t *testing.T) {
	u := assay.Host("https://localhost")
	if u.String() != "https://localhost" {
		t.Error("Unable to make default host")
	}

	u = assay.Host("https://localhost:8080")
	if u.String() != "https://localhost:8080" {
		t.Error("Unable to make default host")
	}

	u = assay.Host("https://localhost:8080/")
	if u.String() != "https://localhost:8080/" {
		t.Error("Unable to make default host")
	}

	u = assay.Host("")
	if u.String() != "" {
		t.Error("Unable to make default host")
	}
}

func TestHostBuildEnv(t *testing.T) {
	for _, env := range []string{
		"https://localhost",
		"https://localhost:8080",
		"https://localhost.localdomain",
		"https://localhost.localdomain:8080",
		"https://localhost.localdomain:8080/",
	} {
		os.Setenv("BUILD_ENDPOINT", env)
		if host := assay.Host(""); host.String() != env {
			t.Errorf("Unexpected host environment %s", host)
		}
	}
	os.Unsetenv("BUILD_ENDPOINT")
}

func TestHostBuildComposition(t *testing.T) {
	os.Setenv("BUILD_ID", "0")
	os.Setenv("CONFIG_ENDPOINT", "https://v%s.example.com")

	if host := assay.Host(""); host.String() != "https://v0.example.com" {
		t.Errorf("Unexpected host composition %s", host)
	}
	os.Unsetenv("BUILD_ID")
	os.Unsetenv("CONFIG_ENDPOINT")
}
