package service

import (
	"net/url"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
var timeoutSecond = 5

type Game8MonsterScoreScrapingResult struct {
	Name        string
	LeaderPoint string
	SubPoint    string
	AssistPoint string
}

type Game8MonsterScrapingResult struct {
	// Noには何かしらのそのページのモンスターに紐づくNoが入るはず
	No     int
	URL    string
	Scores []*Game8MonsterScoreScrapingResult
}

func (r Game8MonsterScrapingResult) ToEntity() entity.Game8Monster {
	scores := make([]entity.Game8MonsterScore, 0, len(r.Scores))
	for _, score := range r.Scores {
		scores = append(scores, entity.Game8MonsterScore{
			Name:        score.Name,
			LeaderPoint: score.LeaderPoint,
			SubPoint:    score.SubPoint,
			AssistPoint: score.AssistPoint,
		})
	}

	url, _ := url.Parse(r.URL)

	return entity.Game8Monster{
		No:     r.No,
		URL:    url.String(),
		Scores: scores,
	}
}

type jsonMonsterData struct {
	No        int    `json:"no"`
	Name      string `json:"name"`
	Evolution struct {
		BaseNo *int `json:"baseNo"`
	} `json:"evolution"`
}

// type gamewithMonsterScrapingResult struct {
// }

// 全モンスターをlocal dynamodbからメモリに展開
func LoadAllMonsters(url *url.URL, tableName string) []*entity.Monster {
	return nil
}
