package service

import (
	"net/url"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

// モンスターを作成
func InitializeMonster(monsterSourceData entity.MonsterSourceData) entity.Monster {
	return entity.Monster{
		No:                 monsterSourceData.No,
		Name:               monsterSourceData.Name,
		OriginMonsterNo:    0,
		Game8URL:           &url.URL{},
		Game8Scores:        []entity.Game8MonsterScore{},
		BatchProcessStatus: 0,
	}
}
