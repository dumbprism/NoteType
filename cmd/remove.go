package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

//remove files that are existing
func removeFile(filename string){
	err := os.Remove(filename + ".md")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(filename + " has been removed")
}
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the specified file",
	Args: cobra.ExactArgs(1),
	Long: ` The remove command is used for removing entries that are present and you have written.`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileName = args[0]
		removeFile(fileName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
