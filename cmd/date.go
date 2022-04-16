/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

// dateCmd represents the date command
var dateCmd = &cobra.Command{
	Use:   "date",
	Short: "格式化时间戳",
	Long: `格式化时间戳

如不传时区参数-z="+8"，默认时区+8，返回格式为 2021-04-16 22:09:33
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("输入错误")
			return
		}
		t, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("请输入时间戳")
			return
		}
		h, _ := cmd.Flags().GetInt("timezone")
		zone := time.FixedZone("CST", 3600*h)
		fmt.Println(time.Unix(t, 0).In(zone).Format("2006-01-02 15:04:05"))
	},
}

func init() {
	rootCmd.AddCommand(dateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dateCmd.Flags().IntP("timezone", "z", +8, "时区 例如 +0800")
}
