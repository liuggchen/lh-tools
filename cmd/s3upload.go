/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

var s3upload = &cobra.Command{
	Use:   "s3upload",
	Short: "上传文件到AWS S3",
	Long: `上传文件到AWS S3

lh-tool s3upload -u accessId -p secretKey -b bucketName -d . -t testTarget

-b 存储桶名称
-d 要上传的目录或文件
-t 目标目录
`,
	Run: func(cmd *cobra.Command, args []string) {
		accessId, _ := cmd.Flags().GetString("accessId")
		secretKey, _ := cmd.Flags().GetString("secretKey")
		region, _ := cmd.Flags().GetString("region")
		bucketName, _ := cmd.Flags().GetString("bucketName")

		dir, _ := cmd.Flags().GetString("dir")
		targetDir, _ := cmd.Flags().GetString("targetDir")

		if dir == "" {
			log.Fatalf("dir is required")
		}

		newSession, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
			Region:      aws.String(region),
		})
		if err != nil {
			log.Fatalln("newSession error: ", err)
		}

		// 创建S3服务客户端
		svc := s3.New(newSession)

		// 遍历目录中的文件并上传
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 跳过目录，只处理文件
			if info.IsDir() {
				return nil
			}

			relativePath, _ := filepath.Rel(dir, path)

			var targetPath string
			if relativePath == "." {
				targetPath = filepath.ToSlash(filepath.Join(targetDir, path))
			} else {
				targetPath = filepath.ToSlash(filepath.Join(targetDir, relativePath))
			}

			if err := uploadFileToS3(svc, path, bucketName, targetPath); err != nil {
				fmt.Printf("error: %s\n", path)
				return err
			} else {
				fmt.Printf("success: %s\n", path)
			}

			return nil
		})

		if err != nil {
			log.Fatalf("failed to walk directory: %v", err)
		}
	},
}

func uploadFileToS3(svc *s3.S3, filePath string, bucketName string, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	fileBytes, _ := io.ReadAll(file)

	_, _ = file.Seek(0, 0)

	contentType := http.DetectContentType(fileBytes)
	if strings.HasPrefix(contentType, "text/plain") {
		ext := filepath.Ext(filePath)
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			contentType = mimeType
		}
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ACL:         aws.String("public-read"), // 可选，设置访问控制
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %q to bucket %q, %v", filePath, bucketName, err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(s3upload)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	s3upload.Flags().StringP("accessId", "u", "", "AWS accessId")
	s3upload.Flags().StringP("secretKey", "p", "", "AWS secretKey")
	s3upload.Flags().StringP("region", "r", "us-east-1", "Region")
	s3upload.Flags().StringP("bucketName", "b", "", "Bucket name")

	s3upload.Flags().StringP("dir", "d", "", "上传目录")
	s3upload.Flags().StringP("targetDir", "t", "", "目标目录")
}
