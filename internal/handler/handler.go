package handler

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *dynamodb.Client
}

func NewHandler(db *dynamodb.Client) *Handler {
	return &Handler{db: db}
}

// GetHelloHandler は、GET /hello のリクエストを処理するハンドラーです。
func (h *Handler) HelloHandler(c *gin.Context) {
	c.HTML(200, "hello.html", gin.H{
		"message": "hello, world",
	})
}

func (h *Handler) GetMonsterByNoHandler(c *gin.Context) {
	monsterNo := c.Param("monster_no")
	c.HTML(200, "monster.html", gin.H{
		"No": monsterNo,
	})
}
