package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)


func createAndAddFile(filename string, title string, entry string, newLineContent string) {

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

	var fullEntry = entry
	
	if newLineContent != "" {
		fullEntry = entry + "\n" + newLineContent
	}

	file.WriteString(structure + "\n" + fullEntry)

	// slice to store all files in the slice

	var allFiles[]string
	allFiles  = append(allFiles,filename)
	fmt.Println("file " + filename + " has been successfully created.")
	fmt.Println("The list of files are  : ")
	for i:=0;i<len(allFiles);i++{
		fmt.Println(allFiles[i])
	}
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
		
		var newLineEntry,err = cmd.Flags().GetString("newline")
		if err != nil{
			fmt.Println(err)
		}
		
		createAndAddFile(filename, title,entry,newLineEntry)
	},
}

func init() {
	newCmd.Flags().StringP("newline","n","","helps to add content in new line")
	rootCmd.AddCommand(newCmd)

}
