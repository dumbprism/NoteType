package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"sort"
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
	tagsView
	templatesView
	themesView
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
	Edit     key.Binding
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
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
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
func (m menuItem) String() string      { return m.Title() + "\n  " + m.Description() }

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
func (n noteItem) String() string      { return n.Title() + "\n  " + n.Description() }

// Tag item
type tagItem struct {
	tag   string
	count int
}

func (t tagItem) Title() string       { return "🏷️  #" + t.tag }
func (t tagItem) Description() string { return fmt.Sprintf("%d entries", t.count) }
func (t tagItem) FilterValue() string { return t.tag }
func (t tagItem) String() string      { return t.Title() + "\n  " + t.Description() }

// Template item
type templateItem struct {
	name string
	desc string
}

func (t templateItem) Title() string       { return "📋 " + t.name }
func (t templateItem) Description() string { return t.desc }
func (t templateItem) FilterValue() string { return t.name }
func (t templateItem) String() string      { return t.Title() + "\n  " + t.Description() }

// Theme item
type themeItem struct {
	name    string
	display string
	current bool
}

func (t themeItem) Title() string {
	if t.current {
		return "✓ 🎨 " + t.display
	}
	return "  🎨 " + t.display
}
func (t themeItem) Description() string { return "Press Enter to apply" }
func (t themeItem) FilterValue() string { return t.name }
func (t themeItem) String() string {
	if t.current {
		return "✓ " + t.display + " (Current)"
	}
	return t.display
}

// Model
type model struct {
	mode         viewMode
	width        int
	height       int
	menuList     list.Model
	notesList    list.Model
	journalsList list.Model
	tagsList     list.Model
	templatesList list.Model
	themesList   list.Model
	editor       textarea.Model
	viewer       viewport.Model
	statusMsg    string
	currentNote  string
	isJournal    bool
	showHelp     bool
	selectedMenu int
}

// Custom delegate for themed list items
type themedDelegate struct {
	list.DefaultDelegate
}

func (d themedDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	str := item.(fmt.Stringer).String()
	
	fn := normalItemStyle.Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}
	
	fmt.Fprint(w, fn(str))
}

func newThemedDelegate() themedDelegate {
	d := themedDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
	}
	return d
}

func initialTUIModel() model {
	// Load and apply theme
	theme := loadTheme()
	applyThemeToStyles(theme)
	
	// Menu items
	items := []list.Item{
		menuItem{title: "Today's Journal", desc: "Write or view today's journal entry", icon: "📔"},
		menuItem{title: "All Journals", desc: "Browse all your journal entries", icon: "📚"},
		menuItem{title: "Notes", desc: "Manage your notes", icon: "📝"},
		menuItem{title: "New Note", desc: "Create a new note", icon: "✨"},
		menuItem{title: "Templates", desc: "Create from template", icon: "📋"},
		menuItem{title: "Tags", desc: "Browse notes by tags", icon: "🏷️"},
		menuItem{title: "Search", desc: "Search across all entries", icon: "🔍"},
		menuItem{title: "Themes", desc: "Change TUI appearance", icon: "🎨"},
		menuItem{title: "Export", desc: "Export to PDF/HTML", icon: "📤"},
		menuItem{title: "Settings", desc: "Configure NoteType", icon: "⚙️"},
	}

	menuList := list.New(items, newThemedDelegate(), 0, 0)
	menuList.Title = "NoteType - Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.Styles.Title = titleStyle
	menuList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)

	// Initialize textarea for editor
	ta := textarea.New()
	ta.Placeholder = "Start writing your thoughts..."
	ta.Focus()
	ta.CharLimit = 0
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(bgAltColor)
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(textColor)

	// Initialize viewport for viewer
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor)

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

		if m.mode == listView || m.mode == tagsView || m.mode == templatesView || m.mode == themesView {
			m.notesList.SetSize(msg.Width-4, msg.Height-8)
			m.journalsList.SetSize(msg.Width-4, msg.Height-8)
			m.tagsList.SetSize(msg.Width-4, msg.Height-8)
			m.templatesList.SetSize(msg.Width-4, msg.Height-8)
			m.themesList.SetSize(msg.Width-4, msg.Height-8)
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
			switch {
			case key.Matches(msg, keys.Edit):
				return m.editCurrentNote()
			default:
				m.viewer, cmd = m.viewer.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case tagsView:
			switch {
			case key.Matches(msg, keys.Enter):
				selectedItem := m.tagsList.SelectedItem()
				if item, ok := selectedItem.(tagItem); ok {
					return m.showEntriesWithTag(item.tag)
				}
			default:
				m.tagsList, cmd = m.tagsList.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case templatesView:
			switch {
			case key.Matches(msg, keys.Enter):
				selectedItem := m.templatesList.SelectedItem()
				if item, ok := selectedItem.(templateItem); ok {
					return m.createFromTemplate(item.name)
				}
			default:
				m.templatesList, cmd = m.templatesList.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case themesView:
			switch {
			case key.Matches(msg, keys.Enter):
				selectedItem := m.themesList.SelectedItem()
				if item, ok := selectedItem.(themeItem); ok {
					return m.applyTheme(item.name)
				}
			default:
				m.themesList, cmd = m.themesList.Update(msg)
				cmds = append(cmds, cmd)
			}
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
	case tagsView:
		content = m.tagsList.View()
	case templatesView:
		content = m.templatesList.View()
	case themesView:
		content = m.themesList.View()
	}

	// Status bar
	status := m.renderStatusBar()

	// Help text
	help := m.renderHelp()

	// Combine all elements with background
	page := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		status,
		help,
	)
	
	// Apply full background color
	return lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Width(m.width).
		Height(m.height).
		Render(page)
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
		Background(bgColor).
		Bold(true).
		MarginBottom(1).
		Render(headerText)

	editorBox := editorStyle.Width(m.width - 4).Render(m.editor.View())

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		activeButtonStyle.Render("💾 Save (Ctrl+S)"),
		inactiveButtonStyle.Render("❌ Cancel (Esc)"),
	)

	return lipgloss.NewStyle().
		Background(bgColor).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			editorBox,
			"\n",
			buttons,
		))
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
		Background(bgColor).
		Bold(true).
		MarginBottom(1).
		Render("👁️  Viewing: " + m.currentNote + " (Press 'e' to edit)")

	viewerBox := panelStyle.Width(m.width - 4).Render(m.viewer.View())

	return lipgloss.NewStyle().
		Background(bgColor).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			viewerBox,
		))
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
	case tagsView:
		modeStr = "🏷️  Tags"
	case templatesView:
		modeStr = "📋 Templates"
	case themesView:
		modeStr = "🎨 Themes"
	}

	left := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(modeStr + " • ") +
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
  
  Actions:       n             New entry (in lists)
                 d             Delete (in lists)
                 e             Edit (in viewer)
                 /             Search
                 Ctrl+S        Save (in editor)
                 ?             Toggle help
  
  TUI Features:
  • Tags: Select from menu to browse all tags
  • Templates: Select to create from template
  • Themes: Select to change colors instantly
  
  Press ? again to hide help
  `

	return lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor).
		Width(m.width - 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
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
	case "Templates":
		return m.loadTemplates()
	case "Tags":
		return m.loadTags()
	case "Search":
		m.mode = searchView
		m.statusMsg = "Search feature"
		return m, nil
	case "Themes":
		return m.loadThemes()
	case "Export":
		m.statusMsg = "Export: Use CLI - notetype export <file>"
		return m, nil
	case "Settings":
		m.statusMsg = "Settings not yet implemented"
		return m, nil
	}
	return m, nil
}

// Load tags view
func (m model) loadTags() (tea.Model, tea.Cmd) {
	tagCounts, err := getAllTags()
	if err != nil {
		m.statusMsg = "Error loading tags: " + err.Error()
		return m, nil
	}
	
	if len(tagCounts) == 0 {
		m.statusMsg = "No tags found. Add #tags to your notes!"
		return m, nil
	}
	
	// Sort by count
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
	
	// Create list items
	var items []list.Item
	for _, tc := range tags {
		items = append(items, tagItem{
			tag:   tc.tag,
			count: tc.count,
		})
	}
	
	m.tagsList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.tagsList.Title = "🏷️  All Tags - Press Enter to filter"
	m.tagsList.Styles.Title = titleStyle
	m.tagsList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
	m.mode = tagsView
	m.statusMsg = fmt.Sprintf("Found %d tags", len(items))
	
	return m, nil
}

// Show entries with specific tag
func (m model) showEntriesWithTag(tag string) (tea.Model, tea.Cmd) {
	files, err := findFilesByTag(tag)
	if err != nil {
		m.statusMsg = "Error finding files: " + err.Error()
		return m, nil
	}
	
	if len(files) == 0 {
		m.statusMsg = fmt.Sprintf("No entries found with #%s", tag)
		return m, nil
	}
	
	// Create list items
	var items []list.Item
	for _, file := range files {
		info, _ := os.Stat(file)
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, ".md")
		items = append(items, noteItem{
			filename: name,
			title:    name,
			date:     info.ModTime().Format("Jan 2, 2006 15:04"),
			size:     formatSizeInTUI(info.Size()),
		})
	}
	
	m.notesList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.notesList.Title = fmt.Sprintf("📄 Entries tagged with #%s", tag)
	m.notesList.Styles.Title = titleStyle
	m.notesList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
	m.mode = listView
	m.isJournal = false
	m.statusMsg = fmt.Sprintf("Found %d entries with #%s", len(items), tag)
	
	return m, nil
}

// Load templates view
func (m model) loadTemplates() (tea.Model, tea.Cmd) {
	templates := []string{"daily", "meeting", "project", "weekly", "idea", "grateful"}
	
	var items []list.Item
	for _, name := range templates {
		items = append(items, templateItem{
			name: name,
			desc: getTemplateDescription(name),
		})
	}
	
	m.templatesList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.templatesList.Title = "📋 Templates - Press Enter to use"
	m.templatesList.Styles.Title = titleStyle
	m.templatesList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
	m.mode = templatesView
	m.statusMsg = fmt.Sprintf("%d templates available", len(items))
	
	return m, nil
}

// Create from template
func (m model) createFromTemplate(templateName string) (tea.Model, tea.Cmd) {
	// Get template content
	templateContent, exists := builtInTemplates[templateName]
	if !exists {
		m.statusMsg = "Template not found"
		return m, nil
	}
	
	// Prepare variables
	now := time.Now()
	vars := map[string]string{
		"date":     now.Format("2006-01-02"),
		"datetime": now.Format("2006-01-02 15:04"),
		"time":     now.Format("15:04"),
		"title":    "New Entry",
		"year":     now.Format("2006"),
		"month":    now.Format("January"),
		"day":      now.Format("Monday"),
	}
	
	// Substitute variables
	finalContent := substituteVariables(templateContent, vars)
	
	// Switch to editor with template content
	m.mode = editorView
	m.isJournal = false
	m.currentNote = fmt.Sprintf("%s-%d", templateName, time.Now().Unix())
	m.editor.SetValue(finalContent)
	m.statusMsg = fmt.Sprintf("Using %s template - Edit and save with Ctrl+S", templateName)
	
	return m, textarea.Blink
}

// Load themes view
func (m model) loadThemes() (tea.Model, tea.Cmd) {
	currentTheme := loadTheme()
	themeNames := []string{"violet", "dracula", "nord", "gruvbox", "solarized", "monokai", "tokyo", "catppuccin"}
	
	var items []list.Item
	for _, name := range themeNames {
		theme := themes[name]
		items = append(items, themeItem{
			name:    name,
			display: theme.Name,
			current: theme.Name == currentTheme.Name,
		})
	}
	
	m.themesList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.themesList.Title = "🎨 Themes - Press Enter to apply"
	m.themesList.Styles.Title = titleStyle
	m.themesList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
	m.mode = themesView
	m.statusMsg = "Select a theme and press Enter"
	
	return m, nil
}

// Apply theme
func (m model) applyTheme(themeName string) (tea.Model, tea.Cmd) {
	theme, exists := themes[themeName]
	if !exists {
		m.statusMsg = "Theme not found"
		return m, nil
	}
	
	// Save theme
	if err := saveTheme(themeName); err != nil {
		m.statusMsg = "Error saving theme: " + err.Error()
		return m, nil
	}
	
	// Apply theme styles
	applyThemeToStyles(theme)
	
	// Recreate menu list with new themed delegate
	items := []list.Item{
		menuItem{title: "Today's Journal", desc: "Write or view today's journal entry", icon: "📔"},
		menuItem{title: "All Journals", desc: "Browse all your journal entries", icon: "📚"},
		menuItem{title: "Notes", desc: "Manage your notes", icon: "📝"},
		menuItem{title: "New Note", desc: "Create a new note", icon: "✨"},
		menuItem{title: "Templates", desc: "Create from template", icon: "📋"},
		menuItem{title: "Tags", desc: "Browse notes by tags", icon: "🏷️"},
		menuItem{title: "Search", desc: "Search across all entries", icon: "🔍"},
		menuItem{title: "Themes", desc: "Change TUI appearance", icon: "🎨"},
		menuItem{title: "Export", desc: "Export to PDF/HTML", icon: "📤"},
		menuItem{title: "Settings", desc: "Configure NoteType", icon: "⚙️"},
	}
	
	m.menuList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.menuList.Title = "NoteType - Main Menu"
	m.menuList.SetShowStatusBar(false)
	m.menuList.SetFilteringEnabled(false)
	m.menuList.Styles.Title = titleStyle
	m.menuList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
	
	// Update editor and viewer styles
	m.editor.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(bgAltColor)
	m.editor.FocusedStyle.Base = lipgloss.NewStyle().Foreground(textColor)
	m.viewer.Style = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor)
	
	m.statusMsg = fmt.Sprintf("✅ Applied theme: %s - All UI elements updated!", theme.Name)
	
	// Go back to menu to see the change
	m.mode = menuView
	
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

	m.journalsList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.journalsList.Title = "📚 Journal Entries"
	m.journalsList.Styles.Title = titleStyle
	m.journalsList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
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

	m.notesList = list.New(items, newThemedDelegate(), m.width-4, m.height-8)
	m.notesList.Title = "📝 Notes"
	m.notesList.Styles.Title = titleStyle
	m.notesList.Styles.TitleBar = lipgloss.NewStyle().
		Background(bgColor).
		Foreground(textColor).
		Padding(0, 1)
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
	filePath := filepath.Join(journalDir, filename+".md")

	content, err := os.ReadFile(filePath)
	if err != nil {
		m.statusMsg = "Error opening journal: " + err.Error()
		return m, nil
	}

	m.mode = viewerView
	m.currentNote = filename
	m.isJournal = true
	m.viewer.SetContent(string(content))
	m.statusMsg = "Viewing journal entry - Press 'e' to edit"
	return m, nil
}

func (m model) openNote(filename string) (tea.Model, tea.Cmd) {
	filePath := filename + ".md"

	content, err := os.ReadFile(filePath)
	if err != nil {
		m.statusMsg = "Error opening note: " + err.Error()
		return m, nil
	}

	m.mode = viewerView
	m.currentNote = filename
	m.isJournal = false
	m.viewer.SetContent(string(content))
	m.statusMsg = "Viewing note - Press 'e' to edit"
	return m, nil
}

func (m model) editCurrentNote() (tea.Model, tea.Cmd) {
	// Load current content into editor
	var filePath string
	if m.isJournal {
		journalDir := getJournalDir()
		filePath = filepath.Join(journalDir, m.currentNote+".md")
	} else {
		filePath = m.currentNote + ".md"
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		m.statusMsg = "Error loading file for editing: " + err.Error()
		return m, nil
	}

	// Switch to editor mode
	m.mode = editorView
	m.editor.SetValue(string(content))
	m.statusMsg = "Editing - Press Ctrl+S to save, Esc to cancel"
	return m, textarea.Blink
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
		filePath := filepath.Join(journalDir, filename+".md")

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			m.statusMsg = "Error saving journal: " + err.Error()
			return m, nil
		}

		m.statusMsg = "✅ Journal saved successfully! Press Esc to go back"
	} else {
		// Save regular note
		filename := m.currentNote
		if filename == "" {
			filename = fmt.Sprintf("note-%d", time.Now().Unix())
		}
		filePath := filename + ".md"

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			m.statusMsg = "Error saving note: " + err.Error()
			return m, nil
		}

		m.statusMsg = "✅ Note saved successfully! Press Esc to go back"
	}

	return m, nil
}

func (m model) deleteSelected() (tea.Model, tea.Cmd) {
	var filePath string

	if m.isJournal {
		if item, ok := m.journalsList.SelectedItem().(noteItem); ok {
			filePath = filepath.Join(getJournalDir(), item.filename+".md")
		}
	} else {
		if item, ok := m.notesList.SelectedItem().(noteItem); ok {
			filePath = item.filename + ".md"
		}
	}

	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
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