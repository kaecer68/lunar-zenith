package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadParityFixturesFromJSON(t *testing.T) {
	tmpDir := t.TempDir()
	p := filepath.Join(tmpDir, "fixtures.json")
	content := `[
  {
    "Name": "sample-a",
    "Date": "2026-04-10",
    "ExpectedOfficer": "開",
    "ExpectedSuitable": ["出行"],
    "ExpectedAvoidable": ["安葬"],
    "ExpectedClash": "沖猴",
    "ExpectedSha": "煞北"
  }
]`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture file failed: %v", err)
	}

	got, err := loadParityFixturesFromJSON(p)
	if err != nil {
		t.Fatalf("loadParityFixturesFromJSON failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("fixture count mismatch: got %d want 1", len(got))
	}
	if got[0].Name != "sample-a" || got[0].ExpectedOfficer != "開" {
		t.Fatalf("fixture content mismatch: %+v", got[0])
	}
}

func TestLoadAllParityFixtures_WithExtraPath(t *testing.T) {
	tmpDir := t.TempDir()
	p := filepath.Join(tmpDir, "fixtures.json")
	content := `[
  {
    "Name": "sample-b",
    "Date": "2026-04-11",
    "ExpectedSuitable": ["出行"]
  }
]`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture file failed: %v", err)
	}

	got, err := loadAllParityFixtures(p)
	if err != nil {
		t.Fatalf("loadAllParityFixtures failed: %v", err)
	}
	if len(got) < 2 {
		t.Fatalf("expected merged fixtures, got len=%d", len(got))
	}
}
