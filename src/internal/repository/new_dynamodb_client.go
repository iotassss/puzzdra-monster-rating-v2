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
	isLocal := os.Getenv("AWS_SAM_LOCAL") == "true"
	isMonsterBatch := os.Getenv("MONSTER_BATCH") == "true"
	if isLocal {
		// ローカルDynamoDB Local用の設定
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					var url string
					if isMonsterBatch {
						url = "http://localhost:8000"
					} else {
						url = "http://host.docker.internal:8000"
					}
					return aws.Endpoint{
						URL:           url,
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
