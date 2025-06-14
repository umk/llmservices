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
	VectorBufSize int

	// Size of an audio buffer used by default.
	AudioBufSize int
}

var Cur = Config{
	RepackPercent: 10,
	VectorBufSize: 1 << 12,
	AudioBufSize:  1 << 21,
}

func Init() error {
	// Define command-line flags
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "Usage: llmservices option...\n")
		fmt.Fprint(w, "Options:\n")
		flag.PrintDefaults()
	}

	flag.IntVar(&Cur.RepackPercent, "repack", Cur.RepackPercent, "percentage of deleted items that triggers repack")
	flag.IntVar(&Cur.VectorBufSize, "vectorbuf", Cur.VectorBufSize, "vector size in vectors pool")
	flag.IntVar(&Cur.AudioBufSize, "audiobuf", Cur.AudioBufSize, "size of audio buffer in bytes")

	// Parse the flags
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Validate repack percentage
	if Cur.RepackPercent <= 0 || Cur.RepackPercent > 100 {
		return errors.New("repack percentage must be between 1 and 100")
	}

	// Validate vector size
	if Cur.VectorBufSize <= 0 {
		return errors.New("vector size must be greater than 0")
	}

	// Validate audio buffer size
	if Cur.AudioBufSize < 100_000 {
		return errors.New("audio buffer size must be at least 100,000 bytes")
	}

	return nil
}
