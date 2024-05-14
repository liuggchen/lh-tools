/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/spf13/cobra"
)

var cfinvalidation = &cobra.Command{
	Use:   "cfinvalidation",
	Short: "添加CloudFront失效",
	Long: `添加CloudFront失效

lh-tool cfinvalidation -u accessId -p secretKey -r "us-east-1" -d distributionId -t "/*,/nouse.txt"

`,
	Run: func(cmd *cobra.Command, args []string) {
		accessId, _ := cmd.Flags().GetString("accessId")
		secretKey, _ := cmd.Flags().GetString("secretKey")
		region, _ := cmd.Flags().GetString("region")
		distributionId, _ := cmd.Flags().GetString("distributionId")

		invalidationPaths, _ := cmd.Flags().GetStringSlice("invalidationPaths")

		var paths []*string
		for _, item := range invalidationPaths {
			if item == "" {
				continue
			}
			paths = append(paths, aws.String(item))
		}

		if len(paths) == 0 {
			log.Fatalln("invalidationPaths is required")
		}

		newSession, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
			Region:      aws.String(region),
		})
		if err != nil {
			log.Fatalln("newSession error: ", err)
		}

		// 创建S3服务客户端
		svc := cloudfront.New(newSession)

		invalidationInput := &cloudfront.CreateInvalidationInput{
			DistributionId: aws.String(distributionId),
			InvalidationBatch: &cloudfront.InvalidationBatch{
				CallerReference: aws.String(fmt.Sprintf("invalidation-%d", time.Now().Unix())),
				Paths: &cloudfront.Paths{
					Quantity: aws.Int64(int64(len(paths))),
					Items:    paths,
				},
			},
		}

		// 创建失效请求
		result, err := svc.CreateInvalidation(invalidationInput)
		if err != nil {
			log.Fatalf("failed to create invalidation, %v", err)
		}

		fmt.Printf("Invalidation created successfully: %s\n", *result.Invalidation.Id)
	},
}

func init() {
	rootCmd.AddCommand(cfinvalidation)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	cfinvalidation.Flags().StringP("accessId", "u", "", "AWS accessId")
	cfinvalidation.Flags().StringP("secretKey", "p", "", "AWS secretKey")
	cfinvalidation.Flags().StringP("region", "r", "us-east-1", "Region")
	cfinvalidation.Flags().StringP("distributionId", "d", "", "Distribution ID ")

	cfinvalidation.Flags().StringSliceP("invalidationPaths", "t", nil, "目标目录 /abc/*")
}
