package auth

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"
)

type Manager struct {
	mu     sync.RWMutex
	tokens map[string]struct{}
	file   *os.File
}

func NewManager(filePath string) (*Manager, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &Manager{
		tokens: make(map[string]struct{}),
		file:   f,
	}, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.file.Close()
}

func (m *Manager) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (m *Manager) AddToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("empty token")
	}

	m.tokens[token] = struct{}{}

	return m.addToFile(token)
}

func (m *Manager) addToFile(token string) error {
	_, err := m.file.WriteString(token + "\n")
	return err
}

func (m *Manager) ValidateToken(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.tokens[strings.TrimSpace(token)]
	return ok
}

func (m *Manager) LoadFromFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(m.file)
	for scanner.Scan() {
		tok := strings.TrimSpace(scanner.Text())
		if tok != "" {
			if _, exists := m.tokens[tok]; !exists {
				m.tokens[tok] = struct{}{}
			}
		}
	}

	return scanner.Err()
}
