package messageHandlers

import (
	"Cerberus/config"
	"testing"
)

func Test_generateTable(t *testing.T) {
	type Service struct {
		Name    string `yaml:"name"`
		Service string `yaml:"service"`
		Ports   string `yaml:"ports"`
	}

	conf := config.ServiceConfig{
		Services: []struct {
			Name    string `yaml:"name"`
			Service string `yaml:"service"`
			Ports   string `yaml:"ports"`
		}{
			{
				Service: "factorio",
				Name:    "Factorio",
			},
			{
				Service: "eco-server",
				Name:    "Eco",
			},
		},
	}

	var resultMap map[string]string = map[string]string{
		"factorio":   "active",
		"eco-server": "failed",
	}

	maxLength := len("Factorio")

	result := generateTable(maxLength, conf, resultMap)
	println(result)
}
