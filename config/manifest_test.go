package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseManifest_Valid(t *testing.T) {
	yamlContent := `
name: test-manifest
containers:
  - name: app
    host: host1
    image: nginx:latest
    entrypoint: /bin/bash
    cmd: -c "echo hello"
    ports:
      - "80:80"
    environment:
      ENV_VAR: value
    options:
      - "--net=my-net"
      - "-v /mnt:/mnt"
`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "manifest.yaml")

	if err := os.WriteFile(file, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp manifest file: %v", err)
	}

	m, err := ParseManifest(file)
	if err != nil {
		t.Fatalf("unexpected error parsing manifest: %v", err)
	}

	if m.Name != "test-manifest" {
		t.Errorf("expected manifest name 'test-manifest', got %q", m.Name)
	}

	if len(m.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(m.Containers))
	}

	c := m.Containers[0]
	if c.Name != "app" {
		t.Errorf("expected container name 'app', got %q", c.Name)
	}
	if c.Host != "host1" {
		t.Errorf("expected container host 'host1', got %q", c.Host)
	}
	if c.Image != "nginx:latest" {
		t.Errorf("expected container image 'nginx:latest', got %q", c.Image)
	}
	if c.Entrypoint != "/bin/bash" {
		t.Errorf("expected entrypoint '/bin/bash', got %q", c.Entrypoint)
	}
	if c.Cmd != `-c "echo hello"` {
		t.Errorf("expected cmd '-c \"echo hello\"', got %q", c.Cmd)
	}
	if len(c.Ports) != 1 || c.Ports[0] != "80:80" {
		t.Errorf("expected ports ['80:80'], got %v", c.Ports)
	}
	if val, ok := c.Environment["ENV_VAR"]; !ok || val != "value" {
		t.Errorf("expected environment ENV_VAR='value', got %v", c.Environment)
	}
	if len(c.Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(c.Options))
	}
}

func TestParseManifest_MissingName(t *testing.T) {
	yamlContent := `
containers:
  - name: app
    host: host1
    image: nginx:latest
`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "manifest.yaml")

	if err := os.WriteFile(file, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp manifest file: %v", err)
	}

	_, err := ParseManifest(file)
	if err == nil {
		t.Fatal("expected error for missing manifest name, got nil")
	}
}

func TestParseManifest_EmptyContainers(t *testing.T) {
	yamlContent := `
name: test-manifest
containers: []
`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "manifest.yaml")

	if err := os.WriteFile(file, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp manifest file: %v", err)
	}

	_, err := ParseManifest(file)
	if err == nil {
		t.Fatal("expected error for empty containers list, got nil")
	}
}

func TestContainerValidation(t *testing.T) {
	tests := []struct {
		name      string
		container Container
		wantErr   bool
	}{
		{
			name: "valid container",
			container: Container{
				Name:  "app",
				Host:  "host1",
				Image: "nginx",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			container: Container{
				Host:  "host1",
				Image: "nginx",
			},
			wantErr: true,
		},
		{
			name: "missing host",
			container: Container{
				Name:  "app",
				Image: "nginx",
			},
			wantErr: true,
		},
		{
			name: "missing image",
			container: Container{
				Name: "app",
				Host: "host1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := tt.container.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestParseManifest_InvalidYAML(t *testing.T) {
	invalidYAML := `
name: test
containers:
  - name: app
    host host1
    image: nginx
`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "invalid.yaml")

	if err := os.WriteFile(file, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	_, err := ParseManifest(file)
	if err == nil {
		t.Fatal("expected YAML parsing error, got nil")
	}
}
