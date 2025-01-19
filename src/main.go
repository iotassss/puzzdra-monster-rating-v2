package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/handler"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/repository"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	h := handler.NewHandler(repository.NewMonsterRepository())

	r.GET("/hello", h.HelloHandler)
	r.GET("/monsters/:no", h.GetMonsterByNoHandler)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginadapter.New(r).Proxy(req)
	})
}
