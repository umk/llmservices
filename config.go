package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/umk/llmservices/internal/config"
	clienthandlers "github.com/umk/llmservices/internal/service/handlers/client"
	"github.com/umk/llmservices/pkg/client"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Clients map[string]client.Config `json:"clients" yaml:"clients" validate:"dive"`
}

func initConfig() error {
	// Initialize the global configuration
	if err := config.Init(); err != nil {
		return err
	}

	if config.Cur.Config == "" {
		return nil
	}

	// Load the configuration from the specified file or STDIN
	c, err := readConfig(config.Cur.Config)
	if err != nil {
		return err
	}

	// Set the clients in the service configuration
	if c.Clients == nil {
		return nil
	}

	for id, conf := range c.Clients {
		c, err := client.New(&conf)
		if err != nil {
			return fmt.Errorf("failed to create client %q: %w", id, err)
		}

		clienthandlers.SetClient(id, c)
	}

	return nil
}

// ReadConfigFromFile reads the configuration from a file path.
func ReadConfigFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	return readConfigFromReader(path, file)
}

// readConfigFromReader reads and parses the config from the given reader and reference name.
func readConfigFromReader(reference string, r io.Reader) (Config, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config

	var unmarshalErr error
	if strings.HasSuffix(reference, ".yaml") || strings.HasSuffix(reference, ".yml") {
		unmarshalErr = yaml.Unmarshal(buf.Bytes(), &config)
	} else {
		unmarshalErr = json.Unmarshal(buf.Bytes(), &config)
	}

	if unmarshalErr != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", unmarshalErr)
	}

	// Validate the configuration using the validator package
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(config); err != nil {
		var valErrors validator.ValidationErrors
		if errors.As(err, &valErrors) && len(valErrors) > 0 {
			e := valErrors[0] // get just the first validation error
			return Config{}, fmt.Errorf("config validation failed: %s: %s", e.Field(), e.Tag())
		}
		return Config{}, fmt.Errorf("config validation error: %w", err)
	}

	return config, nil
}

func readConfig(reference string) (Config, error) {
	r, err := getConfigReader(reference)
	if err != nil {
		return Config{}, err
	}
	defer r.Close()
	return readConfigFromReader(reference, r)
}

func getConfigReader(reference string) (io.ReadCloser, error) {
	if reference == "-" {
		// Read only the first line from standard input when reference is "-"
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return io.NopCloser(strings.NewReader(scanner.Text())), nil
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read from STDIN: %w", err)
		}
		return nil, errors.New("no input provided on STDIN")
	} else {
		// Open the file referenced by the path
		file, err := os.Open(reference)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		return file, nil
	}
}
