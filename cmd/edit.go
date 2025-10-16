package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bernd/desk/internal/shell"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [desk-name]",
	Short: "Edit (or create) a deskfile",
	Long: `Edit (or create) a deskfile with the name specified, otherwise
edit the active deskfile.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) {
	var deskNameToEdit string

	if len(args) == 0 {
		if cfg.DeskName == "" {
			fmt.Fprintln(os.Stderr, "No desk activated.")
			os.Exit(3)
		}
		deskNameToEdit = cfg.DeskName
	} else {
		deskNameToEdit = args[0]
	}

	deskExt, err := shell.GetDeskfileExtension()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to detect shell: %v\n", err)
		os.Exit(1)
	}

	editPath := filepath.Join(cfg.DesksDir, deskNameToEdit+deskExt)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	editorCmd := exec.Command(editor, editPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run editor: %v\n", err)
		os.Exit(1)
	}
}
