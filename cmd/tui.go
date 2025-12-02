package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Styles
var (
	// Color palette
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#8B5CF6")
	accentColor    = lipgloss.Color("#A78BFA")
	successColor   = lipgloss.Color("#10B981")
	warningColor   = lipgloss.Color("#F59E0B")
	errorColor     = lipgloss.Color("#EF4444")
	textColor      = lipgloss.Color("#E5E7EB")
	mutedColor     = lipgloss.Color("#9CA3AF")
	bgColor        = lipgloss.Color("#1F2937")
	bgAltColor     = lipgloss.Color("#374151")

	// Title bar style
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(bgAltColor).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Selected item style
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				PaddingLeft(2)

	// Normal item style
	normalItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	// Panel styles
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

	// Editor styles
	editorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1)

	// Button styles
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
)

// View modes
type viewMode int

const (
	menuView viewMode = iota
	listView
	editorView
	viewerView
	searchView
)

// Key bindings
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Save     key.Binding
	Search   key.Binding
	Delete   key.Binding
	NewEntry key.Binding
	Help     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	NewEntry: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// Menu items
type menuItem struct {
	title string
	desc  string
	icon  string
}

func (m menuItem) Title() string       { return m.icon + " " + m.title }
func (m menuItem) Description() string { return m.desc }
func (m menuItem) FilterValue() string { return m.title }

// Note item
type noteItem struct {
	filename string
	title    string
	date     string
	size     string
}

func (n noteItem) Title() string       { return "📄 " + n.title }
func (n noteItem) Description() string { return n.date + " • " + n.size }
func (n noteItem) FilterValue() string { return n.title }

// Model
type model struct {
	mode         viewMode
	width        int
	height       int
	menuList     list.Model
	notesList    list.Model
	journalsList list.Model
	editor       textarea.Model
	viewer       viewport.Model
	statusMsg    string
	currentNote  string
	isJournal    bool
	showHelp     bool
	selectedMenu int
}

func initialTUIModel() model {
	// Menu items
	items := []list.Item{
		menuItem{title: "Today's Journal", desc: "Write or view today's journal entry", icon: "📔"},
		menuItem{title: "All Journals", desc: "Browse all your journal entries", icon: "📚"},
		menuItem{title: "Notes", desc: "Manage your notes", icon: "📝"},
		menuItem{title: "New Note", desc: "Create a new note", icon: "✨"},
		menuItem{title: "Search", desc: "Search across all entries", icon: "🔍"},
		menuItem{title: "Settings", desc: "Configure NoteType", icon: "⚙️"},
	}

	menuList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menuList.Title = "NoteType - Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.Styles.Title = titleStyle

	// Initialize textarea for editor
	ta := textarea.New()
	ta.Placeholder = "Start writing your thoughts..."
	ta.Focus()
	ta.CharLimit = 0

	// Initialize viewport for viewer
	vp := viewport.New(0, 0)

	return model{
		mode:         menuView,
		menuList:     menuList,
		editor:       ta,
		viewer:       vp,
		statusMsg:    "Welcome to NoteType! Press ? for help",
		selectedMenu: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update component sizes
		m.menuList.SetSize(msg.Width-4, msg.Height-8)
		m.editor.SetWidth(msg.Width - 6)
		m.editor.SetHeight(msg.Height - 12)
		m.viewer.Width = msg.Width - 6
		m.viewer.Height = msg.Height - 12

		if m.mode == listView {
			m.notesList.SetSize(msg.Width-4, msg.Height-8)
			m.journalsList.SetSize(msg.Width-4, msg.Height-8)
		}

	case tea.KeyMsg:
		// Global key bindings
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, keys.Back):
			if m.mode != menuView {
				m.mode = menuView
				m.statusMsg = "Returned to main menu"
				return m, nil
			}
		}

		// Mode-specific key bindings
		switch m.mode {
		case menuView:
			switch {
			case key.Matches(msg, keys.Enter):
				selectedItem := m.menuList.SelectedItem()
				if item, ok := selectedItem.(menuItem); ok {
					return m.handleMenuSelection(item.title)
				}
			default:
				m.menuList, cmd = m.menuList.Update(msg)
				cmds = append(cmds, cmd)
			}

		case editorView:
			switch {
			case key.Matches(msg, keys.Save):
				return m.saveCurrentNote()
			default:
				m.editor, cmd = m.editor.Update(msg)
				cmds = append(cmds, cmd)
			}

		case listView:
			switch {
			case key.Matches(msg, keys.Enter):
				if m.isJournal {
					selectedItem := m.journalsList.SelectedItem()
					if item, ok := selectedItem.(noteItem); ok {
						return m.openJournal(item.filename)
					}
				} else {
					selectedItem := m.notesList.SelectedItem()
					if item, ok := selectedItem.(noteItem); ok {
						return m.openNote(item.filename)
					}
				}
			case key.Matches(msg, keys.NewEntry):
				return m.createNewEntry()
			case key.Matches(msg, keys.Delete):
				return m.deleteSelected()
			default:
				if m.isJournal {
					m.journalsList, cmd = m.journalsList.Update(msg)
				} else {
					m.notesList, cmd = m.notesList.Update(msg)
				}
				cmds = append(cmds, cmd)
			}

		case viewerView:
			m.viewer, cmd = m.viewer.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string

	// Title bar
	title := titleStyle.Width(m.width).Render("✨ NoteType - Your Personal Journal & Notes")

	// Main content based on mode
	switch m.mode {
	case menuView:
		content = m.menuList.View()
	case editorView:
		content = m.renderEditor()
	case listView:
		content = m.renderList()
	case viewerView:
		content = m.renderViewer()
	case searchView:
		content = "Search view (coming soon)"
	}

	// Status bar
	status := m.renderStatusBar()

	// Help text
	help := m.renderHelp()

	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		status,
		help,
	)
}

func (m model) renderEditor() string {
	headerText := "📝 Writing"
	if m.isJournal {
		headerText = "📔 Today's Journal - " + time.Now().Format("Monday, January 2, 2006")
	} else if m.currentNote != "" {
		headerText = "📄 Editing: " + m.currentNote
	}

	header := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		MarginBottom(1).
		Render(headerText)

	editorBox := editorStyle.Width(m.width - 4).Render(m.editor.View())

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		activeButtonStyle.Render("💾 Save (Ctrl+S)"),
		inactiveButtonStyle.Render("❌ Cancel (Esc)"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		editorBox,
		"\n",
		buttons,
	)
}

func (m model) renderList() string {
	if m.isJournal {
		return m.journalsList.View()
	}
	return m.notesList.View()
}

func (m model) renderViewer() string {
	header := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		MarginBottom(1).
		Render("👁️  Viewing: " + m.currentNote)

	viewerBox := panelStyle.Width(m.width - 4).Render(m.viewer.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewerBox,
	)
}

func (m model) renderStatusBar() string {
	// Left side - current mode and status
	var modeStr string
	switch m.mode {
	case menuView:
		modeStr = "📋 Menu"
	case editorView:
		modeStr = "✏️  Editor"
	case listView:
		modeStr = "📚 List"
	case viewerView:
		modeStr = "👁️  Viewer"
	case searchView:
		modeStr = "🔍 Search"
	}

	left := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(modeStr+" • ") +
		lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(m.statusMsg)

	// Right side - time
	right := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(time.Now().Format("15:04"))

	// Create status bar
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 0
	}

	return statusBarStyle.
		Width(m.width).
		Render(left + strings.Repeat(" ", gap) + right)
}

func (m model) renderHelp() string {
	if !m.showHelp {
		return helpStyle.Render("Press ? for help")
	}

	helpText := `
  📌 Keyboard Shortcuts:
  
  Navigation:     ↑/k ↓/j      Move up/down
                 enter         Select/Open
                 esc           Back to menu
                 q / Ctrl+C    Quit
  
  Actions:       n             New entry
                 d             Delete
                 /             Search
                 Ctrl+S        Save
                 ?             Toggle help
  
  Press ? again to hide help
  `

	return helpStyle.
		Width(m.width-4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(1, 2).
		Render(helpText)
}

// formatSize formats file size in human-readable format
func formatSizeInTUI(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Helper functions
func (m model) handleMenuSelection(title string) (tea.Model, tea.Cmd) {
	switch title {
	case "Today's Journal":
		return m.openTodayJournal()
	case "All Journals":
		return m.loadJournals()
	case "Notes":
		return m.loadNotes()
	case "New Note":
		return m.createNewNote()
	case "Search":
		m.mode = searchView
		m.statusMsg = "Search feature"
		return m, nil
	case "Settings":
		m.statusMsg = "Settings not yet implemented"
		return m, nil
	}
	return m, nil
}

func (m model) openTodayJournal() (tea.Model, tea.Cmd) {
	m.mode = editorView
	m.isJournal = true
	m.currentNote = time.Now().Format("2006-01-02")
	m.statusMsg = "Writing today's journal"

	// Load existing content if available
	journalDir := getJournalDir()
	filepath := filepath.Join(journalDir, m.currentNote+".md")

	if content, err := os.ReadFile(filepath); err == nil {
		m.editor.SetValue(string(content))
	} else {
		m.editor.SetValue("")
	}

	return m, textarea.Blink
}

func (m model) loadJournals() (tea.Model, tea.Cmd) {
	journalDir := getJournalDir()
	files, err := filepath.Glob(filepath.Join(journalDir, "*.md"))
	if err != nil {
		m.statusMsg = "Error loading journals: " + err.Error()
		return m, nil
	}

	var items []list.Item
	for _, file := range files {
		info, _ := os.Stat(file)
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		items = append(items, noteItem{
			filename: name,
			title:    name,
			date:     info.ModTime().Format("Jan 2, 2006 15:04"),
			size:     formatSizeInTUI(info.Size()),
		})
	}

	m.journalsList = list.New(items, list.NewDefaultDelegate(), m.width-4, m.height-8)
	m.journalsList.Title = "📚 Journal Entries"
	m.journalsList.Styles.Title = titleStyle
	m.mode = listView
	m.isJournal = true
	m.statusMsg = fmt.Sprintf("Found %d journal entries", len(items))

	return m, nil
}

func (m model) loadNotes() (tea.Model, tea.Cmd) {
	files, err := filepath.Glob("*.md")
	if err != nil {
		m.statusMsg = "Error loading notes: " + err.Error()
		return m, nil
	}

	var items []list.Item
	for _, file := range files {
		info, _ := os.Stat(file)
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		items = append(items, noteItem{
			filename: name,
			title:    name,
			date:     info.ModTime().Format("Jan 2, 2006 15:04"),
			size:     formatSizeInTUI(info.Size()),
		})
	}

	m.notesList = list.New(items, list.NewDefaultDelegate(), m.width-4, m.height-8)
	m.notesList.Title = "📝 Notes"
	m.notesList.Styles.Title = titleStyle
	m.mode = listView
	m.isJournal = false
	m.statusMsg = fmt.Sprintf("Found %d notes", len(items))

	return m, nil
}

func (m model) createNewNote() (tea.Model, tea.Cmd) {
	m.mode = editorView
	m.isJournal = false
	m.currentNote = ""
	m.editor.SetValue("")
	m.statusMsg = "Creating new note"
	return m, textarea.Blink
}

func (m model) createNewEntry() (tea.Model, tea.Cmd) {
	if m.isJournal {
		return m.openTodayJournal()
	}
	return m.createNewNote()
}

func (m model) openJournal(filename string) (tea.Model, tea.Cmd) {
	journalDir := getJournalDir()
	filepath := filepath.Join(journalDir, filename+".md")

	content, err := os.ReadFile(filepath)
	if err != nil {
		m.statusMsg = "Error opening journal: " + err.Error()
		return m, nil
	}

	m.mode = viewerView
	m.currentNote = filename
	m.viewer.SetContent(string(content))
	m.statusMsg = "Viewing journal entry"
	return m, nil
}

func (m model) openNote(filename string) (tea.Model, tea.Cmd) {
	filepath := filename + ".md"

	content, err := os.ReadFile(filepath)
	if err != nil {
		m.statusMsg = "Error opening note: " + err.Error()
		return m, nil
	}

	m.mode = viewerView
	m.currentNote = filename
	m.viewer.SetContent(string(content))
	m.statusMsg = "Viewing note"
	return m, nil
}

func (m model) saveCurrentNote() (tea.Model, tea.Cmd) {
	content := m.editor.Value()

	if m.isJournal {
		// Save to journal directory
		if err := ensureJournalDir(); err != nil {
			m.statusMsg = "Error: " + err.Error()
			return m, nil
		}

		journalDir := getJournalDir()
		filename := m.currentNote
		if filename == "" {
			filename = time.Now().Format("2006-01-02")
		}
		filepath := filepath.Join(journalDir, filename+".md")

		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			m.statusMsg = "Error saving journal: " + err.Error()
			return m, nil
		}

		m.statusMsg = "✅ Journal saved successfully!"
	} else {
		// Save regular note
		filename := m.currentNote
		if filename == "" {
			filename = fmt.Sprintf("note-%d", time.Now().Unix())
		}
		filepath := filename + ".md"

		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			m.statusMsg = "Error saving note: " + err.Error()
			return m, nil
		}

		m.statusMsg = "✅ Note saved successfully!"
	}

	return m, nil
}

func (m model) deleteSelected() (tea.Model, tea.Cmd) {
	var filename string

	if m.isJournal {
		if item, ok := m.journalsList.SelectedItem().(noteItem); ok {
			filename = filepath.Join(getJournalDir(), item.filename+".md")
		}
	} else {
		if item, ok := m.notesList.SelectedItem().(noteItem); ok {
			filename = item.filename + ".md"
		}
	}

	if filename != "" {
		if err := os.Remove(filename); err != nil {
			m.statusMsg = "Error deleting: " + err.Error()
		} else {
			m.statusMsg = "✅ Deleted successfully"
			// Reload the list
			if m.isJournal {
				return m.loadJournals()
			}
			return m.loadNotes()
		}
	}

	return m, nil
}

// TUI command (kept for backwards compatibility, but TUI is now default)
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI interface (default behavior)",
	Long: `Launch NoteType's beautiful terminal user interface.

NOTE: The TUI now launches by default when you run 'notetype' without arguments.
This command is kept for backwards compatibility.

The TUI provides an interactive way to manage your journals and notes
with a modern, keyboard-driven interface.

Features:
  • Beautiful, modern interface
  • Quick access to today's journal
  • Browse all journals and notes
  • Full-screen editor
  • Keyboard shortcuts
  • Real-time saving

To use CLI commands instead, use them directly:
  notetype journal "text"
  notetype new <filename> <title>
  notetype list
  etc.
`,
	Run: func(cmd *cobra.Command, args []string) {
		launchTUI()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
