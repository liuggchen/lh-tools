/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "指定root目录和端口开启server",
	Long: `指定root目录和端口开启server

-d root目录
-p 端口`,
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := cmd.Flags().GetString("dir")
		port, _ := cmd.Flags().GetInt("port")
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.Dir(root)))
		adr := fmt.Sprintf(":%d", port)
		fmt.Printf("server running on: http://127.0.0.1%s\n", adr)
		go openBrowser(adr)
		err := http.ListenAndServe(adr, mux)
		if err != nil {
			fmt.Println("start server error: ", err)
		}
	},
}

func openBrowser(adr string) {
	time.Sleep(2 * time.Second)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "start", "http://127.0.0.1"+adr)
		_ = cmd.Start()
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().StringP("dir", "d", ".", "root目录，默认当前目录")
	serverCmd.Flags().IntP("port", "p", 8080, "开启端口，默认8080")
}
