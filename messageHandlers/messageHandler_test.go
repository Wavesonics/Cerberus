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

	result := generateTable(maxLength, conf, resultMap, 1)
	println(result)
}

func Test_calcMemoryUsage(t *testing.T) {
	cmdParsed := "       16317324      546392"
	result := calcMemUsageFromFreeCommand(cmdParsed)
	expected := (546392.0/16317324.0) * 100

	if result != int(expected) {
		t.Errorf("Expected %d to equal %d", result, int(expected))
	}

	result = calcMemUsageFromFreeCommand(" 0 0")

	if result != 0 {
		t.Errorf("Expected %d to equal %d", result, 0)
	}
}
