package handler

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMonsterByNoHandler(c *gin.Context) {
	reqNo := c.Param("no")

	ctx := c.Request.Context()
	tableName := tableNameMonsterRatings
	result, err := h.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"No": &types.AttributeValueMemberS{Value: reqNo},
		},
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if result.Item == nil {
		c.JSON(404, gin.H{"error": "monster not found"})
		return
	}

	var outputScores []gin.H
	for _, scores := range result.Item["Game8Scores"].(*types.AttributeValueMemberL).Value {
		score := scores.(*types.AttributeValueMemberM).Value
		outputScores = append(outputScores, gin.H{
			"Name":        score["Name"].(*types.AttributeValueMemberS).Value,
			"LeaderPoint": score["LeaderPoint"].(*types.AttributeValueMemberS).Value,
			"SubPoint":    score["SubPoint"].(*types.AttributeValueMemberS).Value,
			"AssistPoint": score["AssistPoint"].(*types.AttributeValueMemberS).Value,
		})
	}
	output := gin.H{
		"No":   result.Item["No"].(*types.AttributeValueMemberS).Value,
		"Name": result.Item["Name"].(*types.AttributeValueMemberS).Value,
		"Game8Monster": gin.H{
			"Scores": outputScores,
			"URL":    result.Item["Game8URL"].(*types.AttributeValueMemberS).Value,
		},
	}

	c.HTML(200, "monster.html", output)
}
