/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "比较两个文件的差集或交集",
	Long: `
比较两个文件的差集或交集，文件按行分隔进行比较。

默认显示在第一个文件中存在但在第二个文件中不存在的行（差集）。
使用 --intersect 参数显示两个文件的交集。

示例:
  lh diff file1.txt file2.txt                    # 显示差集
  lh diff --intersect file1.txt file2.txt       # 显示交集
  lh diff -i file1.txt file2.txt                # 显示交集（简写）
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		file2 := args[1]

		// 获取交集标志
		intersect, _ := cmd.Flags().GetBool("intersect")

		// 读取第一个文件的所有行
		lines1, err := readFileLines(file1)
		if err != nil {
			fmt.Printf("读取文件 %s 失败: %v\n", file1, err)
			return
		}

		// 读取第二个文件的所有行
		lines2, err := readFileLines(file2)
		if err != nil {
			fmt.Printf("读取文件 %s 失败: %v\n", file2, err)
			return
		}

		// 将第二个文件的行存储到map中以便快速查找
		lines2Map := make(map[string]bool)
		for _, line := range lines2 {
			lines2Map[line] = true
		}

		if intersect {
			// 计算交集：找出在两个文件中都存在的行，保持第一个文件的顺序
			var intersection []string
			for _, line := range lines1 {
				if lines2Map[line] {
					intersection = append(intersection, line)
				}
			}

			// 输出交集结果
			if len(intersection) == 0 {
				fmt.Printf("文件 %s 和文件 %s 没有共同的行\n", file1, file2)
			} else {
				fmt.Printf("文件 %s 和文件 %s 的交集（按第一个文件顺序）:\n", file1, file2)
				for _, line := range intersection {
					fmt.Println(line)
				}
			}
		} else {
			// 计算差集：找出在第一个文件中存在但在第二个文件中不存在的行，保持原始顺序
			var diff []string
			for _, line := range lines1 {
				if !lines2Map[line] {
					diff = append(diff, line)
				}
			}

			// 输出差集结果
			if len(diff) == 0 {
				fmt.Printf("文件 %s 中的所有行都在文件 %s 中存在\n", file1, file2)
			} else {
				fmt.Printf("在文件 %s 中存在但在文件 %s 中不存在的行（按原文件顺序）:\n", file1, file2)
				for _, line := range diff {
					fmt.Println(line)
				}
			}
		}
	},
}

// 读取文件的所有行
func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func init() {
	rootCmd.AddCommand(diffCmd)

	// 添加交集标志
	diffCmd.Flags().BoolP("intersect", "i", false, "显示两个文件的交集而不是差集")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
