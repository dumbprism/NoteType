package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)


func createAndAddFile(filename string, title string, entry string) {

	file, err := os.Create(filename + ".md")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// writing inside the file
	var currentDate = time.Now().String()[0:10]
	fmt.Println()
	var style_open = `<span style="opacity:0.5">`
	var style_close = "</span>"
	var structure = "# " + title + "\n" + style_open + currentDate + style_close + "\n" + "---"

	file.WriteString(structure + "\n" + entry)

	fmt.Println("file " + filename + " has been successfully created.")
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "helps you to add new entry to your daily journal",
	Args:  cobra.ExactArgs(3),
	Long: `
		This is where you start typing your thoughts and other things that you wish to type down.
		Don't stop and let your thoughts flow. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var filename = args[0]
		var title = args[1]
		var entry = args[2]

		createAndAddFile(filename, title, entry)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}
