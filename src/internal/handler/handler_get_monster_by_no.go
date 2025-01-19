package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMonsterByNoHandler(c *gin.Context) {
	reqNo := c.Param("no")
	no, err := strconv.Atoi(reqNo)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	ctx := c.Request.Context()

	monster, err := h.monsterRepo.FindByNo(ctx, no)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if monster.No == 0 {
		c.JSON(404, gin.H{"error": "monster not found"})
		return
	}

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

	c.HTML(200, "monster.html", output)
}
