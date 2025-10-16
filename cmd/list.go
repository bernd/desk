package cmd

import (
	"fmt"
	"os"

	"github.com/bernd/desk/internal/deskfile"
	"github.com/bernd/desk/internal/output"
	"github.com/spf13/cobra"
)

var (
	onlyNames bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all desks along with a description",
	Run:     runList,
}

func init() {
	listCmd.Flags().BoolVar(&onlyNames, "only-names", false, "List only the names of the desks")
	listCmd.Flags().BoolVar(&noFormat, "no-format", false, "Use ' - ' to separate names from descriptions")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(cfg.DesksDir); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No desk dir! Run 'desk init'.")
		os.Exit(1)
	}

	desks, err := deskfile.List(cfg.DesksDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing desks: %v\n", err)
		os.Exit(1)
	}

	printDeskList(desks, onlyNames, noFormat)
}

func printDeskList(desks []*deskfile.Deskfile, onlyNames, noFormat bool) {
	if onlyNames {
		for _, desk := range desks {
			fmt.Println(desk.Name)
		}
		return
	}

	var lines []string
	for _, desk := range desks {
		lines = append(lines, fmt.Sprintf("%s - %s", desk.Name, desk.Description))
	}

	if noFormat {
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {
		longest := output.FindLongest(lines)
		output.PrintAligned(lines, longest)
	}
}
