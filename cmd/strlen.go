/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// strlenCmd represents the strlen command
var strlenCmd = &cobra.Command{
	Use:   "strlen",
	Short: "返回字符串长度",
	Long:  `返回字符串长度`,
	Run: func(cmd *cobra.Command, args []string) {
		str := ""
		for _, arg := range args {
			str = fmt.Sprintf("%s %d", str, len(arg))
		}
		fmt.Println("输入的字符串长度为：", str)
	},
}

func init() {
	rootCmd.AddCommand(strlenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// strlenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// strlenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
