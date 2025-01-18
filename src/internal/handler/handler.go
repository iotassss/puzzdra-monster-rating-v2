package handler

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var tableNameMonsters string

type Handler struct {
	db *dynamodb.Client
}

func NewHandler(db *dynamodb.Client) *Handler {
	// とりあえずここで環境変数を読み込む
	if v := os.Getenv("DYNAMODB_TABLE_NAME"); v != "" {
		tableNameMonsters = v
	}

	return &Handler{db: db}
}
