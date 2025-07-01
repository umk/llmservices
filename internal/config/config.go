package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/umk/llmservices/internal/service/handlers"
	"github.com/umk/llmservices/pkg/client"
)

type Config struct {
	// Unix domain socket path. If set, serve from this socket instead of stdio.
	Socket string
	// Path to a configuration file. If not specified, a default file is used.
	File string
	// Name of the default client to override the config file's Default.
	Default string
}

var Cur = Config{
	Socket:  "",
	File:    "",
	Default: "",
}

func Init() error {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: llmservices option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&Cur.Socket, "socket", Cur.Socket, "unix domain socket path to serve from instead of stdio")
	flag.StringVar(&Cur.File, "config", Cur.File, "path to a configuration file")
	flag.StringVar(&Cur.Default, "default", Cur.Default, "ID of default client")

	// Parse the flags
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	f, err := readConfigFiles()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := initClients(f); err != nil {
		return fmt.Errorf("failed to initialize clients: %w", err)
	}

	return nil
}

func initClients(config ConfigFile) error {
	if config.Clients == nil {
		return nil
	}

	clients := make(map[string]*client.Client)

	for id, conf := range config.Clients {
		c, err := client.New(&conf)
		if err != nil {
			return err
		}
		clients[id] = c
	}

	for id, c := range clients {
		handlers.SetGlobalClient(id, c)
	}

	// Handle default client logic
	def := Cur.Default
	if def == "" {
		def = config.Default
	}

	if def != "" {
		if c, ok := clients[def]; ok {
			handlers.SetGlobalClient("default", c)
		} else {
			return fmt.Errorf("default client %q not found", def)
		}
	} else if len(clients) == 1 {
		for _, c := range clients {
			handlers.SetGlobalClient("default", c)
			break
		}
	}

	return nil
}
