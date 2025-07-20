/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var pageDetails struct {
	dateTime string
}

// File Creation -> The file format must be markdown
func createAndWriteFile(filename string, title string) {
	// creation part
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Println("file " + filename + " was created succesfully.")
	fmt.Println(filepath.Ext(filename))

	// writing part -> for now it gives a generic template of file creation.

	pageDetails.dateTime = time.Now().String()[0:10]

	var opacity = "opacity:0.5"
	var style_opening = "<span style=" + opacity+ ">"
	var style_closing = "</span>"
	
	_, err = file.WriteString("# " + title + "\n" + style_opening + pageDetails.dateTime + style_closing + "\n" + "---")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File written sucessfully")
	}
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "If you are willing to add a new entry, this is your starting point",
	Long: `
		The new command helps you to make new entries and create new notes
		to get you starting and brainstorm all your ideas
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var filename = args[0]
		var title = args[1]
		createAndWriteFile(filename, title)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// no flags as of now
}
