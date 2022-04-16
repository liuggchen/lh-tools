/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// lowerCmd represents the lower command
var lowerCmd = &cobra.Command{
	Use:   "lower",
	Short: "字符串转小写",
	Long: `
字符串转小写
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			fmt.Println(strings.ToLower(arg))
		}
	},
}

func init() {
	rootCmd.AddCommand(lowerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lowerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lowerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
