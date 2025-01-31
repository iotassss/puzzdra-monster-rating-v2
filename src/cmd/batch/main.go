package main

import (
	"flag"
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

	batch := batch.NewBatch(repository.NewMonsterRepository())

	fetchMonsterJSON := flag.Bool("fetchmonsterjson", false, "fetch monster json")
	collectGame8MonsterURLs := flag.Bool("collectgame8monsterurls", false, "collect game8 monster urls")
	flag.Parse()

	if *fetchMonsterJSON {
		log.Println("#1 fetch monster json")
		if err := batch.FetchMonsterJSON(); err != nil {
			panic(err)
		}
	} else {
		log.Println("#1 skip fetch monster json")
	}
	if *collectGame8MonsterURLs {
		log.Println("#2 collect game8 monster urls")
		if err := batch.CollectGame8MonsterURLs(); err != nil {
			panic(err)
		}
	} else {
		log.Println("#2 skip collect game8 monster urls")
	}
	log.Println("#3 create all monsters")
	if err := batch.CreateAllMonsters(); err != nil {
		panic(err)
	}
	log.Println("#4 set each site data")
	if err := batch.SetSiteData(); err != nil {
		panic(err)
	}

}
