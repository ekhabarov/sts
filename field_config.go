package sts

import (
	"errors"

	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
)

var ErrConfigMapNotFound = errors.New(`"map" node not found in config`)

type FieldConfig struct {
	// If true, match all fields by Name, even if those fields not in map.
	// Default: false.
	AllMatched bool
	// Field connections.
	FieldMap map[string]string
}

func LoadFieldConfigMap(path string) (*FieldConfig, error) {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	if err := config.LoadFiles(path); err != nil {
		return nil, err
	}

	all, _ := config.Bool("all-matched")

	m, _ := config.StringMap("map")

	return &FieldConfig{
		AllMatched: all,
		FieldMap:   m,
	}, nil
}
