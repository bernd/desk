package config

import (
	"os"
	"path/filepath"
)

const (
	Version       = "0.7.0-beta.1"
	DeskfileName  = "Deskfile"
)

type Config struct {
	DeskDir      string
	DesksDir     string
	ShellArgs    string
	DeskName     string
	DeskEnv      string
}

func Load() *Config {
	homeDir, _ := os.UserHomeDir()
	deskDir := os.Getenv("DESK_DIR")
	if deskDir == "" {
		deskDir = filepath.Join(homeDir, ".desk")
	}

	desksDir := os.Getenv("DESK_DESKS_DIR")
	if desksDir == "" {
		desksDir = filepath.Join(deskDir, "desks")
	}

	return &Config{
		DeskDir:   deskDir,
		DesksDir:  desksDir,
		ShellArgs: os.Getenv("DESK_SHELL_ARGS"),
		DeskName:  os.Getenv("DESK_NAME"),
		DeskEnv:   os.Getenv("DESK_ENV"),
	}
}
