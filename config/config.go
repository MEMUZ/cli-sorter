package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Rules map[string][]string `json:"rules"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) BuildRuleMap() map[string]string {
	rules := map[string]string{}

	for category, exts := range c.Rules {
		for _, ext := range exts {
			rules[ext] = category
		}
	}

	return rules
}
