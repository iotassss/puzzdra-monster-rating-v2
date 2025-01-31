package batch

import (
	"os"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/repository"
)

const (
	// バッチ処理ステータス
	BatchProcessStatusMonsterCreated = 1
	BatchProcessStatusGame8Fetched   = 2
)

type Batch struct {
	monsterRepo                       *repository.MonsterRepository
	monsterDataJsonFilePath           string
	game8MonsterURLListFilePath       string
	failedGame8MonsterURLListFilePath string
	monsterSourceDataJsonURL          string
}

func NewBatch(monsterRepo *repository.MonsterRepository) *Batch {
	monsterDataJsonFilePath := os.Getenv("MONSTER_DATA_JSON_FILE_PATH")
	game8MonsterURLListFilePath := os.Getenv("GAME8_MONSTER_URL_LIST_FILE_PATH")
	failedGame8MonsterURLListFilePath := os.Getenv("FAILED_GAME8_MONSTER_URL_LIST_FILE_PATH")
	monsterSourceDataJsonURL := os.Getenv("MONSTER_SOURCE_DATA_JSON_URL")

	return &Batch{
		monsterRepo:                       monsterRepo,
		monsterDataJsonFilePath:           monsterDataJsonFilePath,
		game8MonsterURLListFilePath:       game8MonsterURLListFilePath,
		failedGame8MonsterURLListFilePath: failedGame8MonsterURLListFilePath,
		monsterSourceDataJsonURL:          monsterSourceDataJsonURL,
	}
}
