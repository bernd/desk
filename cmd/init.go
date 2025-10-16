package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bernd/desk/internal/shell"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize desk configuration",
	Run:   runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(cfg.DeskDir); err == nil {
		fmt.Printf("Desk dir already exists at %s\n", cfg.DeskDir)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Where do you want to store your deskfiles? (default: %s): ", cfg.DeskDir)
	newPrefix, _ := reader.ReadString('\n')
	newPrefix = strings.TrimSpace(newPrefix)
	if newPrefix == "" {
		newPrefix = cfg.DeskDir
	}

	if _, err := os.Stat(newPrefix); os.IsNotExist(err) {
		fmt.Printf("%s doesn't exist, attempting to create.\n", newPrefix)
		if err := os.MkdirAll(filepath.Join(newPrefix, "desks"), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
			os.Exit(1)
		}
	}

	shellType, err := shell.GetRunningShell()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect shell: %v\n", err)
		os.Exit(1)
	}

	homeDir, _ := os.UserHomeDir()
	var shellRC string
	switch shellType {
	case "bash":
		shellRC = filepath.Join(homeDir, ".bashrc")
	case "fish":
		shellRC = filepath.Join(homeDir, ".config", "fish", "config.fish")
	case "zsh":
		shellRC = filepath.Join(homeDir, ".zshrc")
	default:
		shellRC = filepath.Join(homeDir, ".bashrc")
	}

	fmt.Printf("Where's your shell rc file? (default: %s): ", shellRC)
	userShellRC, _ := reader.ReadString('\n')
	userShellRC = strings.TrimSpace(userShellRC)
	if userShellRC == "" {
		userShellRC = shellRC
	}

	if _, err := os.Stat(userShellRC); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s doesn't exist\n", userShellRC)
		os.Exit(1)
	}

	f, err := os.OpenFile(userShellRC, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open shell rc file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close shell rc file: %v\n", closeErr)
		}
	}()

	if _, err := f.WriteString("\n# Hook for desk activation\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to shell rc file: %v\n", err)
		os.Exit(1)
	}
	
	var hookLine string
	if shellType == "fish" {
		hookLine = "test -n \"$DESK_ENV\"; and . \"$DESK_ENV\"; or true\n"
	} else {
		hookLine = "[ -n \"$DESK_ENV\" ] && source \"$DESK_ENV\" || true\n"
	}
	
	if _, err := f.WriteString(hookLine); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to shell rc file: %v\n", err)
		os.Exit(1)
	}

	desksPath := filepath.Join(newPrefix, "desks")
	if newPrefix != cfg.DeskDir {
		fmt.Printf("Done. Start adding desks to %s!\n", desksPath)
		fmt.Printf("Don't forget to set DESK_DIR=%s in your shell rc file.\n", newPrefix)
	} else {
		fmt.Printf("Done. Start adding desks to %s!\n", desksPath)
	}
}
