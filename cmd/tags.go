package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// extractTags finds all #tags in content
func extractTags(content string) []string {
	// Match #tag but not ##heading
	re := regexp.MustCompile(`(?:^|[^#\w])#([\w-]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	tagMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			tagMap[strings.ToLower(match[1])] = true
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// getAllTags scans all files and returns tag usage count
func getAllTags() (map[string]int, error) {
	tagCounts := make(map[string]int)

	// Scan journal entries
	journalDir := getJournalDir()
	if _, err := os.Stat(journalDir); err == nil {
		journalFiles, _ := filepath.Glob(filepath.Join(journalDir, "*.md"))
		for _, file := range journalFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			tags := extractTags(string(content))
			for _, tag := range tags {
				tagCounts[tag]++
			}
		}
	}

	// Scan regular notes
	noteFiles, _ := filepath.Glob("*.md")
	for _, file := range noteFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		tags := extractTags(string(content))
		for _, tag := range tags {
			tagCounts[tag]++
		}
	}

	return tagCounts, nil
}

// findFilesByTag returns files containing the specified tag
func findFilesByTag(tag string) ([]string, error) {
	tag = strings.ToLower(tag)
	var matchingFiles []string

	// Search journal entries
	journalDir := getJournalDir()
	if _, err := os.Stat(journalDir); err == nil {
		journalFiles, _ := filepath.Glob(filepath.Join(journalDir, "*.md"))
		for _, file := range journalFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			tags := extractTags(string(content))
			for _, t := range tags {
				if t == tag {
					matchingFiles = append(matchingFiles, file)
					break
				}
			}
		}
	}

	// Search regular notes
	noteFiles, _ := filepath.Glob("*.md")
	for _, file := range noteFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		tags := extractTags(string(content))
		for _, t := range tags {
			if t == tag {
				matchingFiles = append(matchingFiles, file)
				break
			}
		}
	}

	return matchingFiles, nil
}

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage and view tags in your notes",
	Long: `View all tags used across your journals and notes.

Tags are created by using #hashtag syntax in your notes.
For example: "Today I worked on #project #coding"

Examples:
  notetype tags              # List all tags
  notetype tags list         # List all tags with counts
  notetype tags show work    # Show all entries with #work tag
`,
	Run: func(cmd *cobra.Command, args []string) {
		listAllTags()
	},
}

// tagsListCmd lists all tags
var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags with usage counts",
	Run: func(cmd *cobra.Command, args []string) {
		listAllTags()
	},
}

// tagsShowCmd shows entries with a specific tag
var tagsShowCmd = &cobra.Command{
	Use:   "show <tag>",
	Short: "Show all entries with a specific tag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		// Remove # if provided
		tag = strings.TrimPrefix(tag, "#")

		files, err := findFilesByTag(tag)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if len(files) == 0 {
			fmt.Printf("📝 No entries found with tag #%s\n", tag)
			return
		}

		fmt.Printf("\n📌 Found %d entry/entries with #%s:\n\n", len(files), tag)
		for _, file := range files {
			base := filepath.Base(file)
			name := strings.TrimSuffix(base, ".md")
			fmt.Printf("  • %s\n", name)
		}
		fmt.Println()
	},
}

func listAllTags() {
	tagCounts, err := getAllTags()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(tagCounts) == 0 {
		fmt.Println("📝 No tags found. Add tags to your notes using #hashtag syntax")
		return
	}

	// Sort by count (descending)
	type tagCount struct {
		tag   string
		count int
	}
	var tags []tagCount
	for tag, count := range tagCounts {
		tags = append(tags, tagCount{tag, count})
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].count == tags[j].count {
			return tags[i].tag < tags[j].tag
		}
		return tags[i].count > tags[j].count
	})

	fmt.Printf("\n🏷️  All Tags (%d total):\n\n", len(tags))
	for _, tc := range tags {
		fmt.Printf("  #%-20s (%d)\n", tc.tag, tc.count)
	}
	fmt.Println()
	fmt.Println("💡 Use 'notetype tags show <tag>' to see entries with a specific tag")
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
	tagsCmd.AddCommand(tagsShowCmd)
	rootCmd.AddCommand(tagsCmd)
}
