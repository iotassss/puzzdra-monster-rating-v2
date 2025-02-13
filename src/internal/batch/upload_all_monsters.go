package batch

import (
	"context"
)

func (batch *Batch) UploadAllMonsters() error {
	// モンスターを全件取得
	monsters, err := batch.monsterRepo.ScanAll(context.Background())
	if err != nil {
		return err
	}

	// 全モンスター保存(保存先はクラウド上のDynamoDBを想定)
	ctx := context.Background()
	if err = batch.monsterRemoteRepo.SaveAllMonsters(ctx, monsters); err != nil {
		return err
	}

	return nil
}
