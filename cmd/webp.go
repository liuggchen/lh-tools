/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var webpCmd = &cobra.Command{
	Use:   "webp",
	Short: "转换图片到webp",
	Long: `转换图片到webp

-d 扫描文件目录
-q 图片质量`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		quality, _ := cmd.Flags().GetInt("quality")

		exePath, err := os.Executable()
		if err != nil {
			log.Fatalln("current executable error: ", err)
		}
		exeDir := filepath.Dir(exePath)

		exeFile := filepath.Join(exeDir, "libwebp/bin/cwebp.exe")
		log.Println("exec -> ", exeFile)
		_walkPath(dir, exeFile, quality)
	},
}

func _walkPath(dir string, exeFile string, quality int) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("walk path(%s) error: %s\n", path, err.Error())
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileExt := filepath.Ext(info.Name())
		if ext := fileExt; ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			return nil
		}

		fmt.Println("find：", path)
		outFile := strings.TrimSuffix(path, fileExt) + ".webp"
		cmd := exec.Command(exeFile, "-q", strconv.Itoa(quality), path, "-o", outFile)

		if _, err := cmd.Output(); err != nil {
			fmt.Println("cwebp error: ", err)
		} else {
			fmt.Println("output：", outFile)
		}

		return nil
	})
}

func init() {
	rootCmd.AddCommand(webpCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	webpCmd.Flags().IntP("quality", "q", 77, "图片质量，默认77")
	webpCmd.Flags().StringP("dir", "d", ".", "扫描目录，默认当前目录")
}
