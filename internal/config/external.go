package config

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/umk/llmservices/pkg/client"
	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	// A map of global clients available for any session.
	Clients map[string]client.Config `yaml:"clients,omitempty" validate:"dive"`
	// A name of the client that specified as default in a global clients list.
	Default string `yaml:"default,omitempty"`
}

func readConfigFiles() (ConfigFile, error) {
	if Cur.File == "" {
		if p, err := defaultConfigPath(); err == nil {
			return readConfigFile(p, false)
		}
	} else {
		return readConfigFile(Cur.File, true)
	}
	return ConfigFile{}, nil
}

func readConfigFile(path string, required bool) (ConfigFile, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) && !required {
			return ConfigFile{}, nil
		}
		return ConfigFile{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return ConfigFile{}, err
	}

	var f ConfigFile
	if err := yaml.Unmarshal(b, &f); err != nil {
		return ConfigFile{}, err
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(f); err != nil {
		return ConfigFile{}, err
	}

	return f, nil
}

func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "llmservices.yaml"), nil
}
