package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRequiredPortFromEnv(t *testing.T) {
	t.Setenv("LUNAR_REST_PORT", "18080")

	got, err := GetRequiredPort("LUNAR_REST_PORT", "REST_PORT")
	if err != nil {
		t.Errorf("GetRequiredPort returned error: %v", err)
	}
	if got != "18080" {
		t.Errorf("GetRequiredPort = %q; want 18080", got)
	}
}

func TestGetRequiredPortFromEnvFile(t *testing.T) {
	dir := t.TempDir()
	content := "LUNAR_REST_PORT=28080\nREST_PORT=28080\n"
	if err := os.WriteFile(filepath.Join(dir, localPortsFile), []byte(content), 0644); err != nil {
		t.Errorf("write .env.ports: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	if err := os.Chdir(dir); err != nil {
		t.Errorf("chdir: %v", err)
	}

	got, err := GetRequiredPort("LUNAR_REST_PORT", "REST_PORT")
	if err != nil {
		t.Errorf("GetRequiredPort returned error: %v", err)
	}
	if got != "28080" {
		t.Errorf("GetRequiredPort = %q; want 28080", got)
	}
}
