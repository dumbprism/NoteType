package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Theme represents a color theme
type Theme struct {
	Name          string `json:"name"`
	Primary       string `json:"primary"`
	Secondary     string `json:"secondary"`
	Accent        string `json:"accent"`
	Success       string `json:"success"`
	Warning       string `json:"warning"`
	Error         string `json:"error"`
	Text          string `json:"text"`
	Muted         string `json:"muted"`
	Background    string `json:"background"`
	BackgroundAlt string `json:"background_alt"`
}

// Available themes
var themes = map[string]Theme{
	"violet": {
		Name:          "Violet (Default)",
		Primary:       "#7C3AED",
		Secondary:     "#8B5CF6",
		Accent:        "#A78BFA",
		Success:       "#10B981",
		Warning:       "#F59E0B",
		Error:         "#EF4444",
		Text:          "#E5E7EB",
		Muted:         "#9CA3AF",
		Background:    "#1F2937",
		BackgroundAlt: "#374151",
	},
	"dracula": {
		Name:          "Dracula",
		Primary:       "#BD93F9",
		Secondary:     "#FF79C6",
		Accent:        "#8BE9FD",
		Success:       "#50FA7B",
		Warning:       "#F1FA8C",
		Error:         "#FF5555",
		Text:          "#F8F8F2",
		Muted:         "#6272A4",
		Background:    "#282A36",
		BackgroundAlt: "#44475A",
	},
	"nord": {
		Name:          "Nord",
		Primary:       "#5E81AC",
		Secondary:     "#81A1C1",
		Accent:        "#88C0D0",
		Success:       "#A3BE8C",
		Warning:       "#EBCB8B",
		Error:         "#BF616A",
		Text:          "#ECEFF4",
		Muted:         "#4C566A",
		Background:    "#2E3440",
		BackgroundAlt: "#3B4252",
	},
	"gruvbox": {
		Name:          "Gruvbox Dark",
		Primary:       "#B16286",
		Secondary:     "#D3869B",
		Accent:        "#8EC07C",
		Success:       "#B8BB26",
		Warning:       "#FABD2F",
		Error:         "#FB4934",
		Text:          "#EBDBB2",
		Muted:         "#928374",
		Background:    "#282828",
		BackgroundAlt: "#3C3836",
	},
	"solarized": {
		Name:          "Solarized Dark",
		Primary:       "#268BD2",
		Secondary:     "#2AA198",
		Accent:        "#6C71C4",
		Success:       "#859900",
		Warning:       "#B58900",
		Error:         "#DC322F",
		Text:          "#93A1A1",
		Muted:         "#586E75",
		Background:    "#002B36",
		BackgroundAlt: "#073642",
	},
	"monokai": {
		Name:          "Monokai",
		Primary:       "#F92672",
		Secondary:     "#AE81FF",
		Accent:        "#66D9EF",
		Success:       "#A6E22E",
		Warning:       "#E6DB74",
		Error:         "#F92672",
		Text:          "#F8F8F2",
		Muted:         "#75715E",
		Background:    "#272822",
		BackgroundAlt: "#3E3D32",
	},
	"tokyo": {
		Name:          "Tokyo Night",
		Primary:       "#7AA2F7",
		Secondary:     "#BB9AF7",
		Accent:        "#7DCFFF",
		Success:       "#9ECE6A",
		Warning:       "#E0AF68",
		Error:         "#F7768E",
		Text:          "#C0CAF5",
		Muted:         "#565F89",
		Background:    "#1A1B26",
		BackgroundAlt: "#24283B",
	},
	"catppuccin": {
		Name:          "Catppuccin",
		Primary:       "#CBA6F7",
		Secondary:     "#F5C2E7",
		Accent:        "#89DCEB",
		Success:       "#A6E3A1",
		Warning:       "#F9E2AF",
		Error:         "#F38BA8",
		Text:          "#CDD6F4",
		Muted:         "#6C7086",
		Background:    "#1E1E2E",
		BackgroundAlt: "#313244",
	},
}

// getThemeConfigPath returns the path to theme config
func getThemeConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".notetype-theme.json"
	}
	return filepath.Join(home, ".notetype", "theme.json")
}

// loadTheme loads the current theme from config
func loadTheme() Theme {
	configPath := getThemeConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return default theme
		return themes["violet"]
	}

	var themeName string
	if err := json.Unmarshal(data, &themeName); err != nil {
		return themes["violet"]
	}

	if theme, exists := themes[themeName]; exists {
		return theme
	}

	return themes["violet"]
}

// saveTheme saves the current theme to config
func saveTheme(themeName string) error {
	// Ensure directory exists
	configPath := getThemeConfigPath()
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(themeName)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// applyTheme applies a theme to the TUI styles
func applyThemeToStyles(theme Theme) {
	primaryColor = lipgloss.Color(theme.Primary)
	secondaryColor = lipgloss.Color(theme.Secondary)
	accentColor = lipgloss.Color(theme.Accent)
	successColor = lipgloss.Color(theme.Success)
	warningColor = lipgloss.Color(theme.Warning)
	errorColor = lipgloss.Color(theme.Error)
	textColor = lipgloss.Color(theme.Text)
	mutedColor = lipgloss.Color(theme.Muted)
	bgColor = lipgloss.Color(theme.Background)
	bgAltColor = lipgloss.Color(theme.BackgroundAlt)

	// Rebuild styles with new colors
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)

	statusBarStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgAltColor).
		Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginTop(1)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		PaddingLeft(2)

	normalItemStyle = lipgloss.NewStyle().
		Foreground(textColor).
		PaddingLeft(2)

	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		MarginRight(2)

	activePanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		MarginRight(2)

	editorStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1)

	activeButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		Bold(true).
		Padding(0, 3).
		MarginRight(2)

	inactiveButtonStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Background(bgAltColor).
		Padding(0, 3).
		MarginRight(2)
}

// listAvailableThemes shows all available themes
func listAvailableThemes() {
	currentTheme := loadTheme()

	fmt.Println("\n🎨 Available Themes:\n")

	themeNames := []string{"violet", "dracula", "nord", "gruvbox", "solarized", "monokai", "tokyo", "catppuccin"}

	for _, name := range themeNames {
		theme := themes[name]
		indicator := "  "
		if theme.Name == currentTheme.Name {
			indicator = "✓ "
		}
		fmt.Printf("%s%-15s - %s\n", indicator, name, theme.Name)
		fmt.Printf("   Primary: %s, Accent: %s\n", theme.Primary, theme.Accent)
		fmt.Println()
	}

	fmt.Println("💡 Use 'notetype theme set <name>' to change theme")
	fmt.Println("   Then restart the TUI to see changes")
}

// previewTheme shows a preview of a theme
func previewTheme(themeName string) {
	theme, exists := themes[themeName]
	if !exists {
		fmt.Printf("❌ Theme '%s' not found\n", themeName)
		return
	}

	fmt.Printf("\n🎨 Theme Preview: %s\n\n", theme.Name)

	// Color samples
	fmt.Printf("  Primary:    %s ████\n", theme.Primary)
	fmt.Printf("  Secondary:  %s ████\n", theme.Secondary)
	fmt.Printf("  Accent:     %s ████\n", theme.Accent)
	fmt.Printf("  Success:    %s ████\n", theme.Success)
	fmt.Printf("  Warning:    %s ████\n", theme.Warning)
	fmt.Printf("  Error:      %s ████\n", theme.Error)
	fmt.Printf("  Text:       %s ████\n", theme.Text)
	fmt.Printf("  Muted:      %s ████\n", theme.Muted)
	fmt.Printf("  Background: %s ████\n", theme.Background)
	fmt.Println()
}

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage TUI themes",
	Long: `Change the appearance of the TUI with different color themes.

Available themes:
  violet      - Default purple/violet theme
  dracula     - Dark theme with purple and pink
  nord        - Cool blue theme
  gruvbox     - Warm retro theme
  solarized   - Solarized dark
  monokai     - Monokai pro
  tokyo       - Tokyo Night
  catppuccin  - Catppuccin mocha

Examples:
  notetype theme list           # List all themes
  notetype theme set dracula    # Set Dracula theme
  notetype theme preview nord   # Preview Nord theme
`,
	Run: func(cmd *cobra.Command, args []string) {
		currentTheme := loadTheme()
		fmt.Printf("Current theme: %s\n", currentTheme.Name)
		fmt.Println("\nUse 'notetype theme list' to see all available themes")
	},
}

// themeListCmd lists all themes
var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Run: func(cmd *cobra.Command, args []string) {
		listAvailableThemes()
	},
}

// themeSetCmd sets a theme
var themeSetCmd = &cobra.Command{
	Use:   "set <theme-name>",
	Short: "Set the current theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := args[0]

		if _, exists := themes[themeName]; !exists {
			fmt.Printf("❌ Theme '%s' not found\n", themeName)
			fmt.Println("Use 'notetype theme list' to see available themes")
			return
		}

		if err := saveTheme(themeName); err != nil {
			fmt.Printf("❌ Error saving theme: %v\n", err)
			return
		}

		fmt.Printf("✅ Theme set to '%s'\n", themes[themeName].Name)
		fmt.Println("💡 Restart the TUI to see the changes")
	},
}

// themePreviewCmd previews a theme
var themePreviewCmd = &cobra.Command{
	Use:   "preview <theme-name>",
	Short: "Preview a theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		previewTheme(args[0])
	},
}

func init() {
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeSetCmd)
	themeCmd.AddCommand(themePreviewCmd)
	rootCmd.AddCommand(themeCmd)
}
