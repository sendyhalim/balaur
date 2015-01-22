package balaur

import (
	"github.com/golang/glog"
	"github.com/pelletier/go-toml"
)

func loadTomlConfig(path string) *toml.TomlTree {
	conf, err := toml.LoadFile(path)

	if err != nil {
		glog.Fatalf("Error loading config(%s): %s", path, err)
	}

	return conf
}

func fatalIfNil(val interface{}, message string) {
	if val == nil {
		glog.Fatal(message)
	}
}
