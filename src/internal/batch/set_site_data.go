package batch

import (
	"context"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/service"
)

func (batch *Batch) SetSiteData() error {
	// モンスターをメモリに展開
	monsters, err := batch.monsterRepo.ScanAll(context.Background())

	// game8のモンスターのurl一覧を取得
	strGame8MonsterURLs, err := service.LoadGame8MonsterURLs(batch.game8MonsterURLListFilePath)
	if err != nil {
		return err
	}
	game8MonsterURLs := make([]*url.URL, 0, len(strGame8MonsterURLs))
	for _, strURL := range strGame8MonsterURLs {
		u, err := url.Parse(strURL)
		if err != nil {
			return err
		}
		game8MonsterURLs = append(game8MonsterURLs, u)
	}

	var wg sync.WaitGroup

	// game8モンスターの情報を取得
	wg.Add(1)
	go batch.fetchGame8MonsterData(&wg, game8MonsterURLs, monsters)

	// // gamewithモンスターの情報を取得
	// wg.Add(1)
	// go fetchGamewithMonsterData(gamewithMonsterURLs, monsters)

	wg.Wait()

	return nil
}

func (batch *Batch) fetchGame8MonsterData(wg *sync.WaitGroup, urls []*url.URL, monsters []*entity.Monster) error {
	defer wg.Done()

	errorFile, err := os.OpenFile(batch.failedGame8MonsterURLListFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer errorFile.Close()

	monsterMap := make(map[int]*entity.Monster)
	for _, m := range monsters {
		monsterMap[m.No] = m
	}
	// あらかじめ同じ起源モンスターでモンスターを分類しておく（ファミリー）
	monsterFamilyMap := make(map[int][]*entity.Monster)
	for _, m := range monsters {
		monsterFamilyMap[m.OriginMonsterNo] = append(monsterFamilyMap[m.OriginMonsterNo], m)
	}

	for i, url := range urls {
		// urlからモンスターデータを取得
		result, err := service.ScrapeGame8(url)
		if err != nil {
			errorFile.WriteString(url.String() + "\n")
			log.Printf("Failed to scrape %s: %v", url.String(), err)
			time.Sleep(2 * time.Second)
			continue
		}

		// 取得モンスターNoに対応するモンスターをmonsterMapから取得
		monster := monsterMap[result.No]

		// そのモンスターの起源モンスターNoを取得
		originMonsterNo := monster.OriginMonsterNo

		// monsterFamilyMapから起源モンスターNoに対応するモンスターに取得モンスターの情報を適用
		for _, m := range monsterFamilyMap[originMonsterNo] {
			m.Game8URL = url
			m.Game8Scores = result.ToEntity().Scores
			m.BatchProcessStatus = m.BatchProcessStatus | BatchProcessStatusGame8Fetched
		}

		// 保存
		if err := batch.monsterRepo.SaveAllMonsters(context.Background(), monsterFamilyMap[originMonsterNo]); err != nil {
			errorFile.WriteString(url.String() + "\n")
			log.Printf("Failed to save monster %d: %v", result.No, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Game8monster fetched %d/%d, %v", i+1, len(urls), result)
		time.Sleep(2 * time.Second)
	}

	return nil
}
