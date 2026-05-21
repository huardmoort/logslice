// Package config provides runtime configuration for logslice.
//
// It defines the Config struct that carries all settings for a single
// slicing run, a Validate method for consistency checks, and ParseFlags
// for building a Config from command-line arguments.
//
// Typical usage:
//
//	cfg, err := config.ParseFlags(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
package config
