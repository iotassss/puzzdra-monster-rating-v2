package handler

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

func (h *Handler) GetMonsterByNoHandler(c *gin.Context) {
	reqNo := c.Param("no")
	no, err := strconv.Atoi(reqNo)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "invalid monster no"})
		return
	}

	ctx := c.Request.Context()

	monster, err := h.monsterRepo.FindByNo(ctx, no)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if monster.No == 0 {
		c.JSON(404, gin.H{"error": "monster not found"})
		return
	}

	c.HTML(200, "monster.html", presentMonsterRating(monster))
}

func presentMonsterRating(monster entity.Monster) gin.H {
	var outputScores []gin.H
	for _, score := range monster.Game8Scores {
		outputScores = append(outputScores, gin.H{
			"Name":        score.Name,
			"LeaderPoint": score.LeaderPoint,
			"SubPoint":    score.SubPoint,
			"AssistPoint": score.AssistPoint,
		})
	}
	output := gin.H{
		"No":   monster.No,
		"Name": monster.Name,
		"Game8Monster": gin.H{
			"Scores": outputScores,
			"URL":    monster.Game8URL.String(),
		},
	}
	return output
}
