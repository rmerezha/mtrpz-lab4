package config

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
