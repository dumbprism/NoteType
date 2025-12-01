package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// updateFile appends content to an existing file
func updateFile(filename string, content string, interactive bool, addTimestamp bool) error {
	filepath := filename + ".md"

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist. Use 'notetype new' to create it first", filepath)
	}

	// Open file with proper flags and permissions
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var fullContent string

	if interactive {
		// Interactive mode - allow multi-line input
		fmt.Println("\n✍️  Enter your update (press Ctrl+D or type 'EOF' on a new line to finish):")
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

		fullContent = strings.Join(lines, "")
		fmt.Println(strings.Repeat("-", 70))
	} else {
		fullContent = content
	}

	// Add timestamp if requested
	var updateText string
	if addTimestamp {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		updateText = fmt.Sprintf("\n\n---\n**Updated:** %s\n\n%s", timestamp, fullContent)
	} else {
		updateText = "\n\n" + fullContent
	}

	// Append content with proper formatting
	_, err = file.WriteString(updateText)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Printf("\n✅ Successfully updated '%s'\n", filepath)
	return nil
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <filename> [content]",
	Short: "Append content to an existing journal entry",
	Args:  cobra.MinimumNArgs(1),
	Long: `Update command appends new content to an existing journal entry.

This is useful when you want to add more information to a note without
overwriting the existing content. You can use interactive mode for multi-line updates.

Examples:
  # Interactive mode (for multi-paragraph updates)
  notetype update daily-log -I
  
  # Quick single-line update
  notetype update daily-log "Added this thought later in the day"
  
  # Update with timestamp
  notetype update ideas "Another brilliant idea" -t
  
  # Update today's entry
  notetype update $(date +%Y-%m-%d) "Evening reflection"
`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		var content string
		if len(args) > 1 {
			content = args[1]
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		addTimestamp, _ := cmd.Flags().GetBool("timestamp")

		// If no content provided and not interactive, enable interactive mode
		if content == "" && !interactive {
			interactive = true
		}

		if err := updateFile(filename, content, interactive, addTimestamp); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	updateCmd.Flags().BoolP("interactive", "I", false, "Enter interactive mode for multi-line input")
	updateCmd.Flags().BoolP("timestamp", "t", false, "Add timestamp to the update")
	rootCmd.AddCommand(updateCmd)
}
