package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Built-in templates
var builtInTemplates = map[string]string{
	"daily": `# Daily Journal - {{date}}

## Morning 🌅
**Mood:** 
**Energy Level:** /10

**Today's Intentions:**
- 
- 
- 

**Grateful for:**
- 
- 
- 

## Evening 🌙
**Accomplishments:**
- 
- 
- 

**Lessons Learned:**


**Tomorrow's Focus:**
- 
- 
- 

**Rating:** /10

---
#journal #daily
`,

	"meeting": `# Meeting Notes - {{date}}

**Date:** {{datetime}}
**Attendees:** 

## Agenda
1. 
2. 
3. 

## Discussion


## Decisions Made
- 
- 

## Action Items
- [ ] 
- [ ] 
- [ ] 

## Next Steps


---
#meeting #work
`,

	"project": `# Project: {{title}}

**Start Date:** {{date}}
**Status:** Planning

## Overview


## Goals
1. 
2. 
3. 

## Timeline
- **Week 1:** 
- **Week 2:** 
- **Week 3:** 
- **Week 4:** 

## Resources Needed
- 
- 

## Success Metrics


## Notes


---
#project #planning
`,

	"weekly": `# Weekly Review - Week of {{date}}

## 📊 Overview


## ✅ Wins
- 
- 
- 

## 📈 Progress on Goals


## 🤔 Challenges


## 💡 Lessons Learned


## 🎯 Next Week's Focus
1. 
2. 
3. 

---
#weekly-review #reflection
`,

	"idea": `# Idea: {{title}}

**Date:** {{date}}

## The Idea


## Why This Matters


## Next Steps
- [ ] 
- [ ] 
- [ ] 

## Resources


## Notes


---
#ideas #brainstorm
`,

	"grateful": `# Gratitude - {{date}}

Today I'm grateful for:

1. 
2. 
3. 

## Why?


## Reflection


---
#gratitude #reflection
`,
}

// getTemplateDir returns the templates directory path
func getTemplateDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./templates"
	}
	return filepath.Join(home, ".notetype", "templates")
}

// ensureTemplateDir creates the template directory if it doesn't exist
func ensureTemplateDir() error {
	templateDir := getTemplateDir()
	return os.MkdirAll(templateDir, 0755)
}

// substituteVariables replaces template variables with actual values
func substituteVariables(content string, vars map[string]string) string {
	result := content
	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// applyTemplate creates a note from a template
func applyTemplate(templateName, filename, title string) error {
	// Get template content
	var templateContent string
	var exists bool

	// Check built-in templates first
	templateContent, exists = builtInTemplates[templateName]

	// Check custom templates
	if !exists {
		templatePath := filepath.Join(getTemplateDir(), templateName+".md")
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("template '%s' not found", templateName)
		}
		templateContent = string(content)
	}

	// Prepare variables
	now := time.Now()
	vars := map[string]string{
		"date":     now.Format("2006-01-02"),
		"datetime": now.Format("2006-01-02 15:04"),
		"time":     now.Format("15:04"),
		"title":    title,
		"year":     now.Format("2006"),
		"month":    now.Format("January"),
		"day":      now.Format("Monday"),
	}

	// Substitute variables
	finalContent := substituteVariables(templateContent, vars)

	// Create file
	filePath := filename + ".md"
	if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	return nil
}

// listTemplates shows all available templates
func listTemplates() {
	fmt.Println("\n📋 Built-in Templates:\n")

	templates := []string{"daily", "meeting", "project", "weekly", "idea", "grateful"}
	for _, name := range templates {
		desc := getTemplateDescription(name)
		fmt.Printf("  %-15s - %s\n", name, desc)
	}

	// Check for custom templates
	templateDir := getTemplateDir()
	if _, err := os.Stat(templateDir); err == nil {
		customTemplates, _ := filepath.Glob(filepath.Join(templateDir, "*.md"))
		if len(customTemplates) > 0 {
			fmt.Println("\n📝 Custom Templates:\n")
			for _, tmpl := range customTemplates {
				name := strings.TrimSuffix(filepath.Base(tmpl), ".md")
				fmt.Printf("  %s\n", name)
			}
		}
	}

	fmt.Println("\n💡 Usage: notetype template <template-name> <filename> <title>")
	fmt.Println("   Example: notetype template daily today \"My Day\"")
}

func getTemplateDescription(name string) string {
	descriptions := map[string]string{
		"daily":    "Daily journal with morning/evening sections",
		"meeting":  "Meeting notes with agenda and action items",
		"project":  "Project planning template",
		"weekly":   "Weekly review and reflection",
		"idea":     "Capture and develop ideas",
		"grateful": "Gratitude journal entry",
	}
	return descriptions[name]
}

// showTemplate displays a template content
func showTemplate(templateName string) {
	content, exists := builtInTemplates[templateName]
	if !exists {
		templatePath := filepath.Join(getTemplateDir(), templateName+".md")
		contentBytes, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Printf("❌ Template '%s' not found\n", templateName)
			return
		}
		content = string(contentBytes)
	}

	fmt.Printf("\n📄 Template: %s\n", templateName)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(content)
	fmt.Println(strings.Repeat("=", 70))
}

// saveCustomTemplate saves a custom template
func saveCustomTemplate(name, content string) error {
	if err := ensureTemplateDir(); err != nil {
		return err
	}

	templatePath := filepath.Join(getTemplateDir(), name+".md")
	return os.WriteFile(templatePath, []byte(content), 0644)
}

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template <template-name> <filename> <title>",
	Short: "Create notes from templates",
	Long: `Create new notes using pre-defined templates.

Built-in templates:
  daily    - Daily journal template
  meeting  - Meeting notes template
  project  - Project planning template
  weekly   - Weekly review template
  idea     - Idea capture template
  grateful - Gratitude journal template

Examples:
  notetype template daily today "My Daily Entry"
  notetype template meeting standup "Team Standup"
  notetype template project project-x "Project X"
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			// Just show the template
			showTemplate(args[0])
			return
		}

		if len(args) < 3 {
			fmt.Println("❌ Usage: notetype template <template-name> <filename> <title>")
			return
		}

		templateName := args[0]
		filename := args[1]
		title := args[2]

		if err := applyTemplate(templateName, filename, title); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Created '%s.md' from template '%s'\n", filename, templateName)
		fmt.Printf("💡 Edit it with: notetype\n")
	},
}

// templateListCmd lists all templates
var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Run: func(cmd *cobra.Command, args []string) {
		listTemplates()
	},
}

// templateShowCmd shows a template
var templateShowCmd = &cobra.Command{
	Use:   "show <template-name>",
	Short: "Show template content",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		showTemplate(args[0])
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	rootCmd.AddCommand(templateCmd)
}
