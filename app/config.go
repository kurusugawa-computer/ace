package app

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/kurusugawa-computer/ace/app/agent"
)

type Config struct {
	Agents map[string]*agent.Config `yaml:"agents"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.OpenFile(path, 0, os.FileMode(os.O_RDONLY))
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(f, yaml.UseJSONUnmarshaler())
	config := Config{}
	err = dec.Decode(&config)
	f.Close()
	if err != nil {
		return nil, err
	}

	// エージェントConfigのkeyをエージェントのNameとしてセット
	for name, agentConfig := range config.Agents {
		agentConfig.Name = name
	}

	return &config, nil
}
