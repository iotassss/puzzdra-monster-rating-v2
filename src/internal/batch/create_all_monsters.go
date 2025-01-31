package batch

import (
	"context"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/progressbar"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/service"
)

func (batch *Batch) CreateAllMonsters() error {
	// モンスターデータを読み込む
	monsterSourceDataList, err := service.LoadMonsterDataJsonFile(batch.monsterDataJsonFilePath)
	if err != nil {
		return err
	}

	// モンスターを初期化
	monsters := make([]*entity.Monster, 0, len(monsterSourceDataList))
	for _, m := range monsterSourceDataList {
		monster := service.InitializeMonster(*m)
		monsters = append(monsters, &monster)
	}

	// モンスターに対応する起源モンスターを探す
	for i, monster := range monsters {
		originMonsterSourceData, err := service.FindOriginMonsterSourceData(monsterSourceDataList, monster.No)
		if err != nil {
			return err
		}
		monster.OriginMonsterNo = originMonsterSourceData.No

		monster.BatchProcessStatus = monster.BatchProcessStatus | BatchProcessStatusMonsterCreated
		progressbar.Display(i+1, len(monsters))
	}

	// 全モンスター保存
	ctx := context.Background()
	if err = batch.monsterRepo.SaveAllMonsters(ctx, monsters); err != nil {
		return err
	}

	return nil
}
