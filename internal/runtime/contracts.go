package runtime

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const localPortsFile = ".env.ports"

// GetRequiredPort resolves a required runtime port from env vars first and then
// from the generated .env.ports file. It never falls back to a hardcoded port.
func GetRequiredPort(primary string, aliases ...string) (string, error) {
	keys := append([]string{primary}, aliases...)

	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value, nil
		}
	}

	ports, err := loadLocalPorts()
	if err != nil {
		return "", err
	}

	for _, key := range keys {
		if value := strings.TrimSpace(ports[key]); value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("missing required runtime port %s (checked env and %s)", primary, localPortsFile)
}

func loadLocalPorts() (map[string]string, error) {
	path, err := findLocalPortsFile()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	return values, nil
}

func findLocalPortsFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		path := filepath.Join(dir, localPortsFile)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return "", fmt.Errorf("%s not found from working directory %s; run scripts/sync-contracts.sh first", localPortsFile, wd)
}
