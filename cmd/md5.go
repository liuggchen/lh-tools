/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: "返回32位md5值",
	Long: `
返回输入字符串的32位md5值

默认返回小写结果
-u 返回大写结果
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			u, _ := cmd.Flags().GetBool("upper")
			fmt.Println(_md5(arg, u))
		}
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)
	md5Cmd.Flags().BoolP("upper", "u", false, "-u 返回大写")
}

func _md5(str string, isUpper bool) string {
	sum := md5.Sum([]byte(str))
	s := hex.EncodeToString(sum[:])
	if isUpper {
		s = strings.ToUpper(s)
	}
	return s
}
