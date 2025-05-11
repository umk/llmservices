package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	// Vectors database repack percentage threshold
	RepackPercent int

	// Pooled vector size. Can be greater or less to the size of
	// actual vectors operated by a client.
	VectorSize int
}

var C = Config{
	RepackPercent: 10,
	VectorSize:    20_000,
}

func Init() error {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: llmservices option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	flag.IntVar(&C.RepackPercent, "repack", C.RepackPercent, "percentage of deleted items that triggers repack")
	flag.IntVar(&C.VectorSize, "vector", C.VectorSize, "vector size in vectors pool")

	// Parse the flags
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Validate repack percentage
	if C.RepackPercent <= 0 || C.RepackPercent > 100 {
		return errors.New("repack percentage must be between 1 and 100")
	}

	// Validate vector size
	if C.VectorSize <= 0 {
		return errors.New("vector size must be greater than 0")
	}

	return nil
}
