package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// This package contains the configuration for customers service.

func Load() (*Config, error) {
	baseCfg, err := loadFromYAML("config/base.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	overrideFromEnv(baseCfg)

	return baseCfg, nil
}

func loadFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	processStruct(v, "")
}

func processStruct(v reflect.Value, prefix string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		envName := prefix + strings.ToUpper(structField.Name)

		// For nested structs
		if field.Kind() == reflect.Struct {
			processStruct(field, envName+"_")
			continue
		}

		// Get value from env
		envValue, exists := os.LookupEnv(envName)
		if !exists {
			continue
		}

		// Set value
		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			if intVal, err := strconv.Atoi(envValue); err == nil {
				field.SetInt(int64(intVal))
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				sliceVal := strings.Split(envValue, ",")
				field.Set(reflect.ValueOf(sliceVal))
			}
		}
	}
}
