package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func newDynamoDBClient() *dynamodb.Client {
	var cfg aws.Config
	var err error

	// 環境変数でローカル実行かどうかを判定
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		// ローカルDynamoDB Local用の設定
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL: "http://host.docker.internal:8000", // Dockerで起動したDynamoDB Localを指す
						// URL:           "http://localhost:8000", // Dockerで起動したDynamoDB Localを指す
						SigningRegion: "ap-northeast-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
		)
		log.Println("Using DynamoDB Local")
	} else {
		// クラウドDynamoDB用の設定
		cfg, err = config.LoadDefaultConfig(context.TODO())
		log.Println("Using DynamoDB in AWS Cloud")
	}

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}
