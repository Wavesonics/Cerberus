package config

type ServiceConfig struct {
	Service []struct {
		// Human friendly display name
		Name string `yaml:"name"`
		// The literal systemd service name
		Service string `yaml:"service"`
		// A string explaining which ports are used
		Ports string `yaml:"ports"`
	}
}
