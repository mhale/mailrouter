package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Config struct {
	sync.RWMutex
	Routes  map[string]Route
	Filters map[string]Filter
	Options map[string]string
}

// Add the DROP route. Make it the default if there is no existing default route.
func AddDropRoute() {
	dropIsDefault := true
	for _, route := range config.Routes {
		if route.IsDefault == true {
			dropIsDefault = false
		}
	}
	config.Routes["DROP"] = Route{Id: "DROP", Name: "Drop", IsDefault: dropIsDefault}
}

func SetDefaultOptions() {
	if _, exists := config.Options["PIDFile"]; !exists {
		config.Options["PIDFile"] = ""
	}
}

// Load the filter and route configuration from a JSON file.
// Add the drop route as it must always be present.
func LoadConfig() error {
	defer func() {
		if config.Routes == nil {
			config.Routes = map[string]Route{}
		}
		if config.Filters == nil {
			config.Filters = map[string]Filter{}
		}
		if config.Options == nil {
			config.Options = map[string]string{}
		}
		AddDropRoute()
		SetDefaultOptions()
	}()

	data, err := ioutil.ReadFile(*confFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}

// Create a clone of the global config to use in SaveConfig().
func CloneConfig() *Config {
	clone := new(Config)
	clone.Routes = map[string]Route{}
	clone.Filters = map[string]Filter{}
	clone.Options = map[string]string{}
	for k, v := range config.Routes {
		clone.Routes[k] = v
	}
	for k, v := range config.Filters {
		clone.Filters[k] = v
	}
	for k, v := range config.Options {
		clone.Options[k] = v
	}
	return clone
}

// Save the filter and route configuration to a JSON file.
// Remove the DROP route before marshalling - it is a hardcoded route that should never be in the config file.
// Don't delete DROP from the actual config variable (it might be needed mid-save), make a copy instead.
func SaveConfig() error {
	clone := CloneConfig()
	delete(clone.Routes, "DROP")

	data, err := json.MarshalIndent(clone, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(*confFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
