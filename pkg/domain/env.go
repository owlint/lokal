package domain

type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
