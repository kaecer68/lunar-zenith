package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocsConsistencyScriptPasses(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	repoRoot := filepath.Clean(filepath.Join(wd, "..", ".."))
	scriptPath := filepath.Join(repoRoot, "scripts", "check-docs-consistency.sh")

	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check-docs-consistency failed: %v\n%s", err, string(output))
	}
}

func TestWebUIOnlyUsesV4ContractShape(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	repoRoot := filepath.Clean(filepath.Join(wd, "..", ".."))
	uiPath := filepath.Join(repoRoot, "internal", "webui", "static", "index.html")
	content, err := os.ReadFile(uiPath)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", uiPath, err)
	}

	body := string(content)
	forbidden := []string{
		"f?.Name ??",
		"f?.Type ??",
		"f?.Description ??",
		"lunar.Month || lunar.month",
		"lunar.Day || lunar.day",
		"lunar.IsLeap || lunar.is_leap",
		"lunar.Year || lunar.year",
		"pillars.Year || pillars.year",
		"pillars.Month || pillars.month",
		"pillars.Day || pillars.day",
		"pillars.Hour || pillars.hour",
		"StemIndex",
		"BranchIndex",
		"ld.Day",
		"ld.IsLeap",
		"ld.Month",
	}
	for _, needle := range forbidden {
		if strings.Contains(body, needle) {
			t.Errorf("index.html should not contain stale v3 fallback %q", needle)
		}
	}

	required := []string{
		"mansion.full_name",
		"dailyDeity.type",
		"fetalGod.position",
		"clashSha.clash_zodiac",
		"lunar.month",
		"pillars.year",
		"f?.name || ''",
	}
	for _, needle := range required {
		if !strings.Contains(body, needle) {
			t.Errorf("index.html should contain v4 contract anchor %q", needle)
		}
	}
}
