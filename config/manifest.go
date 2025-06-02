package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Manifest struct {
	Name       string      `yaml:"name"`
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name        string            `yaml:"name"`
	Host        string            `yaml:"host"`
	Image       string            `yaml:"image"`
	Entrypoint  string            `yaml:"entrypoint,omitempty"`
	Cmd         string            `yaml:"cmd,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Options     []string          `yaml:"options,omitempty"`
}

func ParseManifest(filename string) (*Manifest, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *Manifest) Validate() error {
	if m.Name == "" {
		return errors.New("manifest: 'name' field is required")
	}
	if len(m.Containers) == 0 {
		return errors.New("manifest: 'containers' field is required and cannot be empty")
	}
	for i, c := range m.Containers {
		if err := c.Validate(); err != nil {
			return errors.New("container[" + c.Name + "]: " + err.Error())
		}
		if c.Name == "" {
			return errors.New("container[" + string(rune(i)) + "]: 'name' field is required")
		}
	}
	return nil
}

func (c *Container) Validate() error {
	if c.Name == "" {
		return errors.New("'name' field is required")
	}
	if c.Host == "" {
		return errors.New("'host' field is required")
	}
	if c.Image == "" {
		return errors.New("'image' field is required")
	}
	return nil
}
