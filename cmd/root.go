package cmd

import (
	"fmt"
	"os"

	"github.com/bernd/desk/internal/config"
	"github.com/bernd/desk/internal/deskfile"
	"github.com/bernd/desk/internal/output"
	"github.com/spf13/cobra"
)

var (
	noFormat bool
	cfg      *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "desk",
	Short: "Lightweight workspace manager for the shell",
	Long:  "◲  desk - Lightweight workspace manager for the shell",
	Run:   runCurrent,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&noFormat, "no-format", false, "Use ' - ' to separate items from descriptions")
}

func initConfig() {
	cfg = config.Load()
}

func runCurrent(cmd *cobra.Command, args []string) {
	if cfg.DeskEnv == "" {
		fmt.Println("No desk activated.")
		fmt.Println()
		
		desks, err := deskfile.List(cfg.DesksDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No desk dir! Run 'desk init'.\n")
			os.Exit(1)
		}
		
		printDeskList(desks, false, noFormat)
		os.Exit(1)
	}

	desk, err := deskfile.Parse(cfg.DeskEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing deskfile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s - %s\n", cfg.DeskName, desk.Description)

	if len(desk.Callables) > 0 {
		var lines []string
		for _, callable := range desk.Callables {
			if callable.Doc != "" {
				lines = append(lines, fmt.Sprintf("  %s - %s", callable.Name, callable.Doc))
			} else {
				lines = append(lines, fmt.Sprintf("  %s", callable.Name))
			}
		}

		fmt.Println()
		if noFormat {
			for _, line := range lines {
				fmt.Println(line)
			}
		} else {
			longest := output.FindLongest(lines)
			output.PrintAligned(lines, longest)
		}
	}
}
