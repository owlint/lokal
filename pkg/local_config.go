package pkg

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ConfigEnvironmentVariable struct {
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Secret string `yaml:"secret"`
}

type Forwarding struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type LocalConfig struct {
	Env     []ConfigEnvironmentVariable `yaml:"env"`
	Forward []Forwarding                `yaml:"forward"`
	Command string                      `yaml:"command"`
}

func ReadLocalConfig(path string) (*LocalConfig, error) {
	configFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}
	config := LocalConfig{}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
