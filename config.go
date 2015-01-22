package balaur

import (
	"fmt"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config interface {
	Get(string, bool) string
	GetArray(string, bool) []string
	GetChildren(string, bool) []Config
}

func NewConfig(path string) Config {
	var c Config

	switch filepath.Ext(path) {
	case ".toml":
		c = NewTomlConfig(path)
	}

	return c
}

func NewTomlConfig(path string) *TomlConfig {
	return &TomlConfig{
		config: loadTomlConfig(path),
	}
}

type TomlConfig struct {
	config *toml.TomlTree
}

func (tmc *TomlConfig) Get(key string, errorIfNotExist bool) string {
	keyInterface := tmc.config.Get(key)
	getValue := func() string {
		return keyInterface.(string)
	}

	if !errorIfNotExist {
		if keyInterface == nil {
			return ""
		} else {
			return getValue()
		}
	}

	tmc.fatalIfNil(keyInterface, key)
	return getValue()
}

func (tmc *TomlConfig) GetArray(key string, errorIfNotExist bool) []string {
	var configs []string
	configInterface := tmc.config.Get(key)

	getValue := func() []string {
		for _, val := range configInterface.([]interface{}) {
			configs = append(configs, val.(string))
		}
		return configs
	}

	if !errorIfNotExist {
		if configInterface == nil {
			return configs
		} else {
			return getValue()
		}
	}

	tmc.fatalIfNil(configInterface, key)

	return getValue()
}

func (tmc *TomlConfig) GetChildren(key string, errorIfNotExist bool) []Config {
	var configs []Config
	configInterface := tmc.config.Get(key)

	getValue := func() []Config {
		for _, val := range configInterface.([]*toml.TomlTree) {
			configs = append(configs, &TomlConfig{
				config: val,
			})
		}
		return configs
	}

	if !errorIfNotExist {
		if configInterface == nil {
			return configs
		} else {
			return getValue()
		}
	}

	tmc.fatalIfNil(configInterface, key)

	return getValue()
}

func (tmc *TomlConfig) fatalIfNil(keyInterface interface{}, key string) {
	fatalIfNil(keyInterface, fmt.Sprintf("Config key(%s) is nil", key))
}
