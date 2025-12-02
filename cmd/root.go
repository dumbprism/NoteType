/*
Copyright © 2025 KOTAMRAJU ARHANT <arhantk915@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "notetype",
	Version: version,
	Short:   "Your one stop destination to create and manage your notes",
	Long: `
NoteType - A beautiful note-taking and journaling application

The sole purpose of NoteType is to give users the feel of CLI and also 
help them journal things out at times when they cannot carry a book around. 
Thus, in this era of digital transformation, it is quite necessary to have it.

By default, NoteType launches a beautiful TUI (Terminal User Interface).
Use CLI commands for scripting and automation.

CLI Commands:
  journal - Daily journaling
  new     - Create a new note
  update  - Append content to an existing note
  remove  - Delete a note
  list    - List all notes
  view    - View the contents of a note
  search  - Search for notes by title or content
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Launch TUI by default
		launchTUI()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("cli", "c", false, "Show CLI help instead of launching TUI")
}

// launchTUI starts the TUI interface
func launchTUI() {
	p := tea.NewProgram(
		initialTUIModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}