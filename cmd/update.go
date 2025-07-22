package cmd

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
)

// function to update files

func updateFile(filename string,content string){
	// appending to the file
	file,err := os.OpenFile(filename+".md",os.O_APPEND,6660)
	
	if err != nil{
		fmt.Println(err)
	}

	if _,err := file.Write([]byte("\n"+content));err != nil{
		file.Close()
		fmt.Println(err)
	}
}


// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Appends data into an existing file",
	Long: `Helps you to update your file if you have missed something or would like to add something
	so that you don't miss out on anything
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var filename = args[0]
		var content = args[1]

		updateFile(filename,content)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.LocalFlags().String("help","noteype update <filename> <content>","Used for updating files")
}
