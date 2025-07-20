package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

func deleteNote(filename string){
	err := os.Remove(filename)
	if err !=nil{
		fmt.Println(err)
	}
	fmt.Println("deleted entry succesfully")
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",  
	Short: "Delete an exisiting file note that you have created",
	Long: `Delete an exisiting file that is present`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filename = args[0]
		deleteNote(filename)
	},
}


func init() {
	rootCmd.AddCommand(deleteCmd)
}
