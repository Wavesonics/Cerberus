package config

type ServiceConfig struct {
	Service []struct {
		Name    string `yaml:"name"`
		Service string `yaml:"service"`
		Ports   string `yaml:"ports"`
	}
}
