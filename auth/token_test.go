package auth_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rmerezha/mtrpz-lab4/auth"
)

func TestManager_TokenLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tokens.txt")

	manager, err := auth.NewManager(filePath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	token, err := manager.GenerateToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	const expectedLength = 64
	if len(token) != expectedLength {
		t.Errorf("token length = %d; want %d", len(token), expectedLength)
	}

	if err := manager.AddToken(token); err != nil {
		t.Fatalf("failed to add token: %v", err)
	}

	if !manager.ValidateToken(token) {
		t.Errorf("expected token to be valid, got invalid")
	}
}

func TestManager_LoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tokens.txt")

	const token = "abcd1234"

	if err := os.WriteFile(filePath, []byte(token+"\n"), 0644); err != nil {
		t.Fatalf("failed to write token to file: %v", err)
	}

	manager, err := auth.NewManager(filePath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	if err := manager.LoadFromFile(); err != nil {
		t.Fatalf("failed to load tokens from file: %v", err)
	}

	if !manager.ValidateToken(token) {
		t.Errorf("expected token to be valid after load, got invalid")
	}
}

func TestManager_AddEmptyToken(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tokens.txt")

	manager, err := auth.NewManager(filePath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	err = manager.AddToken("   ")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestManager_TokenPersists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tokens.txt")

	manager, err := auth.NewManager(filePath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	token := "persist-token-123"
	if err := manager.AddToken(token); err != nil {
		t.Fatalf("add token failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !strings.Contains(string(data), token) {
		t.Errorf("expected token to be persisted in file")
	}
}
