package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Book   BookConfig       `toml:"book"`
	Output map[string]Theme `toml:"output"`
}

type BookConfig struct {
	Title   string   `toml:"title"`
	Authors []string `toml:"authors"`
	Src     string   `toml:"src"`
}

type Theme struct {
	DefaultTheme string `toml:"default-theme"`
	SiteURL      string `toml:"site-url"`
}

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(config *Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
