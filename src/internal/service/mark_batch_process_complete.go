package service

import "github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"

// バッチ処理ステータスに1をもたせる
func MarkBatchProcessComplete(entity.Monster) entity.Monster {
	return entity.Monster{}
}
