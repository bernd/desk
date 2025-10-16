package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/bernd/desk/internal/deskfile"
	"github.com/bernd/desk/internal/shell"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <desk-name> '<cmd>' OR run <desk-name> <cmd> <arg>...",
	Short: "Run a command within a desk's environment then exit",
	Long: `Run a command within a desk's environment then exit. In the first form
shell expansion is performed. In the second form, the argument vector
is executed as is.`,
	Args:               cobra.MinimumNArgs(2),
	Run:                runRun,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) {
	toDesk := args[0]
	cmdArgs := args[1:]

	deskExt, err := shell.GetDeskfileExtension()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect shell: %v\n", err)
		os.Exit(1)
	}

	deskPath, err := deskfile.Find(cfg.DesksDir, toDesk, deskExt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Desk %s (%s%s) not found in %s\n", toDesk, toDesk, deskExt, cfg.DesksDir)
		os.Exit(1)
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

	var execArgs []string
	if len(cmdArgs) == 1 {
		execArgs = []string{shellExec, "-ic", cmdArgs[0]}
	} else {
		quotedArgs := make([]string, len(cmdArgs))
		for i, arg := range cmdArgs {
			quotedArgs[i] = fmt.Sprintf("%q", arg)
		}
		cmdStr := strings.Join(quotedArgs, " ")
		execArgs = []string{shellExec, "-ic", cmdStr}
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("DESK_NAME=%s", toDesk))
	env = append(env, fmt.Sprintf("DESK_ENV=%s", deskPath))

	err = syscall.Exec(shellPath, execArgs, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to exec shell: %v\n", err)
		os.Exit(1)
	}
}
