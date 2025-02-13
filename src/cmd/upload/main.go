package main

import (
	"log"
	"os"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/batch"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/repository"
)

func main() {
	// tmp
	os.Setenv("AWS_SAM_LOCAL", "true")
	os.Setenv("MONSTER_BATCH", "true")
	os.Setenv("DYNAMODB_TABLE_NAME", "LocalMonsters")
	os.Setenv("MONSTER_DATA_JSON_FILE_PATH", "data/monsters.json")
	os.Setenv("GAME8_MONSTER_URL_LIST_FILE_PATH", "data/game8_monster_urls.txt")
	os.Setenv("FAILED_GAME8_MONSTER_URL_LIST_FILE_PATH", "data/failed_game8_monster_urls.txt")
	os.Setenv("MONSTER_SOURCE_DATA_JSON_URL", "https://padmdb.rainbowsite.net/listJson/monster_data.json")
	os.Setenv("DYNAMODB_TABLE_NAME_REMOTE", "PrdMonsters")

	monsterRepo := repository.NewMonsterRepository()
	monsterRemoteRepo := repository.NewMonsterRepositoryWithRemoteDB()
	batch := batch.NewBatch(monsterRepo, monsterRemoteRepo)

	log.Println("#1 upload monsters to remote db")
	if err := batch.UploadAllMonsters(); err != nil {
		panic(err)
	}
}
