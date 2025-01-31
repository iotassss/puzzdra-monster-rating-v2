package service

import (
	"encoding/json"
	"os"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

// jsonファイルから必要なデータを読み込む
func LoadMonsterDataJsonFile(monsterDataJsonFilePath string) ([]*entity.MonsterSourceData, error) {
	// monsterDataJsonFilePathで指定されたファイルを読み込む
	// ファイルの中身を元に全モンスターの元データを作成
	file, err := os.Open(monsterDataJsonFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var sourceMonsters []*entity.MonsterSourceData

	// JSONが配列の形式であるか確認する
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	// JSON配列の各要素を順番にデコードする
	for decoder.More() {
		var jsonMonster jsonMonsterData
		if err := decoder.Decode(&jsonMonster); err != nil {
			return nil, err
		}

		sourceMonster := &entity.MonsterSourceData{
			No:     jsonMonster.No,
			Name:   jsonMonster.Name,
			BaseNo: jsonMonster.Evolution.BaseNo,
		}

		sourceMonsters = append(sourceMonsters, sourceMonster)
	}
	// JSON配列の終了を確認する
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}

	return sourceMonsters, nil
}
