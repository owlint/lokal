package config

import (
	"errors"
	"os"

	"github.com/owlint/lokal/pkg/domain"
	"gopkg.in/yaml.v3"
)

type Forwarding struct {
	FromPort int `yaml:"from-port"`
	ToPort   int `yaml:"to-port"`
}

type LocalConfig struct {
	Env            []domain.EnvironmentVariable `yaml:"env"`
	Forward        []Forwarding                 `yaml:"forward"`
	Command        string                       `yaml:"command"`
	Namespace      string                       `yaml:"namespace"`
	Deployment     string                       `yaml:"deployment"`
	Container      string                       `yaml:"container"`
	ForceNamespace bool                         `yaml:"force-namespace"`
}

func (c LocalConfig) EnsureValid() error {
	if c.Namespace == "" {
		return errors.New("namespace cannot be empty")
	}

	if c.Deployment == "" {
		return errors.New("deployment cannot be empty")
	}

	if c.Container == "" {
		return errors.New("container cannot be empty")
	}

	if c.Command == "" {
		return errors.New("command cannot be empty")
	}

	return nil
}

func ReadLocalConfig(path string) (*LocalConfig, error) {
	configFile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}
	config := LocalConfig{}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
