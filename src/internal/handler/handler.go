package handler

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const tableNameMonsterRatings = "MonsterRatings"

type Handler struct {
	db *dynamodb.Client
}

func NewHandler(db *dynamodb.Client) *Handler {
	return &Handler{db: db}
}
