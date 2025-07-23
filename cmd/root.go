/*
Copyright © 2025 KOTAMRAJU ARHANT <arhantk915@gmail.com>

*/
package cmd

import (
	"os"
	
	"github.com/spf13/cobra"
)
var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:   "notetype",
	Version: version,
	Short: "Your one stop destination to create your add your notes and store them",
	Long: `
		The sole purpose of NoteType is to give users the feel of CLI and also 
		help them journal things out at times when they cannot carry a book around. 
		Thus, in this era of digital transformation, it is quite necessary to have it.
	`,
	Run: func(cmd *cobra.Command, args []string){

	},
	}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


