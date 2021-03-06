package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// The configuration of the CLI
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	Lang struct {
		Default string `yaml:"default"`
	}
	ConfigFileName string `yaml:"config_file_name"`
	Version        string `yaml:"version"`
	Author         string `yaml:"author"`
}

func getDefaultConfigLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "egor.yaml"), nil
}

func createDefaultConfiguration() *Config {
	return &Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 1200,
		},
		Lang: struct {
			Default string `yaml:"default"`
		}{
			Default: "cpp",
		},
		Version:        "0.1.0",
		ConfigFileName: "egor-meta.json",
	}
}

// This function is called when the configuration file does not exist already
// This will create the configuration file in the user config dir, with a minimalistic
// default configuration
func SaveConfiguration(config *Config) error {
	location, err := getDefaultConfigLocation()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, buffer.Bytes(), 0777)
}

// Returns the Configuration object associated with
// the path given as a parameter
func LoadConfiguration(location string) (*Config, error) {
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Returns the Configuration object associated with
// the default configuration location
func LoadDefaultConfiguration() (*Config, error) {
	location, err := getDefaultConfigLocation()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(location); err != nil {
		if os.IsNotExist(err) {
			config := createDefaultConfiguration()
			if err := SaveConfiguration(config); err != nil {
				return nil, err
			}
		}
	}
	return LoadConfiguration(location)
}

// Gets the configuration value associated with the given key
func GetConfigurationValue(config *Config, key string) (string, error) {
	lowerKey := strings.ToLower(key)
	if lowerKey == "server.port" {
		return strconv.Itoa(config.Server.Port), nil
	} else if lowerKey == "lang.default" {
		return config.Lang.Default, nil
	} else if lowerKey == "author" {
		return config.Author, nil
	} else {
		return "", errors.New(fmt.Sprintf("Unknown config key %s", key))
	}
}
