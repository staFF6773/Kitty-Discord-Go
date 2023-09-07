package config

import (
	"encoding/json"
	"os"
)

func ParseConfigFromJSONFile(path string) (*Configuration, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	if err := json.Unmarshal(contents, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
