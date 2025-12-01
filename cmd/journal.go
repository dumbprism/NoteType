package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// getJournalDir returns the journal directory path
func getJournalDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./journal"
	}
	return filepath.Join(home, ".notetype", "journal")
}

// ensureJournalDir creates the journal directory if it doesn't exist
func ensureJournalDir() error {
	journalDir := getJournalDir()
	return os.MkdirAll(journalDir, 0755)
}

// getTodayFilename returns the filename for today's journal entry
func getTodayFilename() string {
	return time.Now().Format("2006-01-02")
}

// createTodayEntry creates or appends to today's journal entry
func createTodayEntry(entry string, interactive bool) error {
	if err := ensureJournalDir(); err != nil {
		return fmt.Errorf("error creating journal directory: %v", err)
	}

	journalDir := getJournalDir()
	filename := getTodayFilename()
	filepath := filepath.Join(journalDir, filename+".md")

	// Check if today's entry already exists
	fileExists := false
	if _, err := os.Stat(filepath); err == nil {
		fileExists = true
	}

	var content string

	if interactive || entry == "" {
		// Interactive mode - allow multi-line input
		fmt.Println("\n📔 Daily Journal Entry")
		fmt.Println(strings.Repeat("=", 70))
		if fileExists {
			fmt.Println("📝 Adding to today's entry...")
		} else {
			fmt.Println("📝 Creating today's entry...")
		}
		fmt.Println("\nWrite your thoughts (press Ctrl+D or type 'EOF' on a new line to finish):")
		fmt.Println(strings.Repeat("-", 70))

		reader := bufio.NewReader(os.Stdin)
		var lines []string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				// EOF reached
				break
			}

			// Check if user typed EOF
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "EOF" || trimmedLine == "eof" {
				break
			}

			lines = append(lines, line)
		}

		content = strings.Join(lines, "")
		fmt.Println(strings.Repeat("-", 70))
	} else {
		content = entry
	}

	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("no content provided")
	}

	if fileExists {
		// Append to existing file
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		timestamp := time.Now().Format("15:04")
		updateText := fmt.Sprintf("\n\n### %s\n\n%s", timestamp, content)

		_, err = file.WriteString(updateText)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}

		fmt.Printf("\n✅ Added entry to today's journal (%s)\n", filename)
	} else {
		// Create new file
		file, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()

		currentDate := time.Now().Format("Monday, January 2, 2006")
		timestamp := time.Now().Format("15:04")

		structure := fmt.Sprintf("# Daily Journal\n\n## %s\n\n### %s\n\n%s",
			currentDate, timestamp, content)

		_, err = file.WriteString(structure)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}

		fmt.Printf("\n✅ Created today's journal entry (%s)\n", filename)
	}

	fmt.Printf("📍 Location: %s\n", filepath)
	return nil
}

// viewTodayEntry displays today's journal entry
func viewTodayEntry() error {
	if err := ensureJournalDir(); err != nil {
		return fmt.Errorf("error accessing journal directory: %v", err)
	}

	journalDir := getJournalDir()
	filename := getTodayFilename()
	filepath := filepath.Join(journalDir, filename+".md")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("no journal entry for today yet. Create one with 'notetype journal'")
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("  📔 Today's Journal Entry (%s)\n", filename)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
	fmt.Println(string(content))
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("📍 %s\n\n", filepath)

	return nil
}

// listJournalEntries lists all journal entries
func listJournalEntries(limit int) error {
	if err := ensureJournalDir(); err != nil {
		return fmt.Errorf("error accessing journal directory: %v", err)
	}

	journalDir := getJournalDir()

	files, err := filepath.Glob(filepath.Join(journalDir, "*.md"))
	if err != nil {
		return fmt.Errorf("error reading journal entries: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("📝 No journal entries yet. Create your first entry with 'notetype journal'")
		return nil
	}

	// Sort files by name (date) in reverse order (newest first)
	type fileInfo struct {
		path string
		name string
	}

	var fileInfos []fileInfo
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		fileInfos = append(fileInfos, fileInfo{path: file, name: name})
	}

	// Simple reverse sort
	for i := len(fileInfos)/2 - 1; i >= 0; i-- {
		opp := len(fileInfos) - 1 - i
		fileInfos[i], fileInfos[opp] = fileInfos[opp], fileInfos[i]
	}

	displayCount := len(fileInfos)
	if limit > 0 && limit < displayCount {
		displayCount = limit
	}

	fmt.Printf("\n📚 Journal Entries (showing %d of %d):\n\n", displayCount, len(fileInfos))

	for i := 0; i < displayCount; i++ {
		info, err := os.Stat(fileInfos[i].path)
		if err != nil {
			continue
		}

		modTime := info.ModTime().Format("15:04")

		// Parse date for better display
		t, err := time.Parse("2006-01-02", fileInfos[i].name)
		var displayDate string
		if err == nil {
			displayDate = t.Format("Mon, Jan 2, 2006")
		} else {
			displayDate = fileInfos[i].name
		}

		fmt.Printf("  📅 %s (last updated: %s)\n", displayDate, modTime)
	}

	fmt.Printf("\n📍 Journal location: %s\n\n", journalDir)
	return nil
}

var journalCmd = &cobra.Command{
	Use:   "journal [entry]",
	Short: "Quick access to daily journaling",
	Long: `The journal command provides quick access to daily journaling.

All journal entries are automatically stored in ~/.notetype/journal/
with dates as filenames (YYYY-MM-DD.md).

Subcommands:
  (no args)  - Create or append to today's entry (interactive mode)
  add        - Add to today's entry (interactive mode)
  view       - View today's entry
  list       - List all journal entries

Examples:
  # Write today's journal (interactive)
  notetype journal
  
  # Quick one-line entry
  notetype journal "Today was amazing!"
  
  # View today's entry
  notetype journal view
  
  # List all entries
  notetype journal list
`,
	Run: func(cmd *cobra.Command, args []string) {
		var entry string
		if len(args) > 0 {
			entry = args[0]
		}

		if err := createTodayEntry(entry, entry == ""); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	},
}

var journalViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View today's journal entry",
	Run: func(cmd *cobra.Command, args []string) {
		if err := viewTodayEntry(); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	},
}

var journalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all journal entries",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		if err := listJournalEntries(limit); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	journalListCmd.Flags().IntP("limit", "l", 0, "Limit number of entries to display (0 = all)")

	journalCmd.AddCommand(journalViewCmd)
	journalCmd.AddCommand(journalListCmd)
	rootCmd.AddCommand(journalCmd)
}
