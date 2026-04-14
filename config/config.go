package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	AppID             string            `yaml:"appId"`
	Scheme            string            `yaml:"scheme"`
	Addr              string            `yaml:"addr"`
	Inspect           bool              `yaml:"inspect"`
	Authorization     string            `yaml:"authorization"`
	RequestHeaderAdd  map[string]string `yaml:"request_header_add,omitempty"`
	ResponseHeaderAdd map[string]string `yaml:"response_header_add,omitempty"`
}

type Configuration struct {
	Version      string               `yaml:"version"`
	Token        string               `yaml:"token"`
	Region       string               `yaml:"region"`
	ConsoleUI    bool                 `yaml:"console_ui"`
	Applications map[string]AppConfig `yaml:"applications"`
	Meta         map[string]string    `yaml:"meta,omitempty"`
}

func GetConfigPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\agglabs\\loop.yml"
	}
	return "/etc/agglabs/loop.yml"
}

func LoadConfiguration() (*Configuration, error) {
	path := GetConfigPath()
	config := &Configuration{
		Version:      "3",
		Region:       "default",
		ConsoleUI:    true,
		Applications: make(map[string]AppConfig),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	for name, app := range config.Applications {
		if app.Scheme == "" {
			app.Scheme = "http"
		}
		if app.Authorization == "" {
			app.Authorization = "public"
		}
		config.Applications[name] = app
	}

	return config, nil
}

func SaveConfiguration(config *Configuration) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func GetPidFilePath(appName string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join("C:\\agglabs", fmt.Sprintf("%s.pid", appName))
	}
	return filepath.Join("/etc/agglabs", fmt.Sprintf("%s.pid", appName))
}
