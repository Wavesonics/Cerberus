package messageHandlers

import (
	"Cerberus/config"
	"fmt"
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

	result := generateTable(maxLength, conf, resultMap, "01%", "1%")
	println(result)
}

func Test_calcMemoryUsage(t *testing.T) {
	cmdParsed := "       16317324      546392"
	result := calcMemUsageFromFreeCommand(cmdParsed)
	expected := fmt.Sprintf("%03.1f%%", (546392.0/16317324.0) * 100)

	if result != expected {
		t.Errorf("Expected %s to equal %s", result, expected)
	}

	result = calcMemUsageFromFreeCommand(" 0 0")

	if result != "0.0%" {
		t.Errorf("Expected %s to equal %s", result, "0.0%")
	}
}

func Test_calcCpuUsage(t *testing.T) {
	cmdParsed := "77.12"
	result := calcCpuUsageFromMpstatCommand(cmdParsed)
	expected := fmt.Sprintf("%03.1f%%", 22.88)

	if result != expected {
		t.Errorf("Expected %s to equal %s", result, expected)
	}
}
