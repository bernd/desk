package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/bernd/desk/internal/config"
	"github.com/bernd/desk/internal/deskfile"
	"github.com/bernd/desk/internal/shell"
	"github.com/spf13/cobra"
)

var goCmd = &cobra.Command{
	Use:     "go [desk-name-or-path] [shell-args...]",
	Aliases: []string{"."},
	Short:   "Activate a desk",
	Long: `Activate a desk. Extra arguments are passed onto shell. If called with
no arguments, look for a Deskfile in the current directory. If not a
recognized desk, try as a path to directory containing a Deskfile.`,
	Run:                runGo,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
	var toDesk string
	var shellArgs []string

	if len(args) > 0 {
		toDesk = args[0]
		shellArgs = args[1:]
	}

	deskExt, err := shell.GetDeskfileExtension()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect shell: %v\n", err)
		os.Exit(1)
	}

	var deskPath string
	var deskName string

	if toDesk == "" {
		deskPath = filepath.Join(".", config.DeskfileName)
		if _, err := os.Stat(deskPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "No Deskfile found in current directory\n")
			os.Exit(1)
		}
		absPath, _ := filepath.Abs(deskPath)
		deskName = filepath.Base(filepath.Dir(absPath))
		deskPath = absPath
	} else {
		possiblePath, err := deskfile.Find(cfg.DesksDir, toDesk, deskExt)
		if err == nil {
			deskPath = possiblePath
			deskName = toDesk
		} else {
			possibleDeskfileDir := toDesk
			if filepath.Base(toDesk) == config.DeskfileName {
				possibleDeskfileDir = filepath.Dir(toDesk)
			}

			if stat, err := os.Stat(possibleDeskfileDir); err == nil && stat.IsDir() {
				deskfilePath := filepath.Join(possibleDeskfileDir, config.DeskfileName)
				if _, err := os.Stat(deskfilePath); err == nil {
					absPath, _ := filepath.Abs(deskfilePath)
					deskPath = absPath
					deskName = filepath.Base(possibleDeskfileDir)
				} else {
					fmt.Fprintf(os.Stderr, "No Deskfile found in '%s'\n", possibleDeskfileDir)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Desk %s (%s%s) not found in %s\n", toDesk, toDesk, deskExt, cfg.DesksDir)
				os.Exit(1)
			}
		}
	}

	shellExec, err := shell.GetRunningShell()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect shell: %v\n", err)
		os.Exit(1)
	}

	shellPath, err := findShellPath(shellExec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find shell executable: %v\n", err)
		os.Exit(1)
	}

	execArgs := []string{shellExec}
	if cfg.ShellArgs != "" {
		execArgs = append(execArgs, cfg.ShellArgs)
	}
	execArgs = append(execArgs, shellArgs...)

	env := os.Environ()
	env = append(env, fmt.Sprintf("DESK_NAME=%s", deskName))
	env = append(env, fmt.Sprintf("DESK_ENV=%s", deskPath))

	err = syscall.Exec(shellPath, execArgs, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to exec shell: %v\n", err)
		os.Exit(1)
	}
}

func findShellPath(shellName string) (string, error) {
	possiblePaths := []string{
		"/bin/" + shellName,
		"/usr/bin/" + shellName,
		"/usr/local/bin/" + shellName,
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("shell %s not found", shellName)
}
