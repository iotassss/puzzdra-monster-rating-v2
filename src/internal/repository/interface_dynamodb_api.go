package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBAPI interface {
	GetItem(
		ctx context.Context,
		params *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
	BatchWriteItem(
		ctx context.Context,
		params *dynamodb.BatchWriteItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.BatchWriteItemOutput, error)
	Scan(
		ctx context.Context,
		params *dynamodb.ScanInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.ScanOutput, error)
}
