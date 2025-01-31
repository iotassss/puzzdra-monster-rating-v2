package service

import (
	"fmt"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

// モンスターに対応する起源モンスターを探す
func FindOriginMonsterSourceData(allMonsters []*entity.MonsterSourceData, targetNo int) (*entity.MonsterSourceData, error) {
	// Noをキーにマップ化
	monsterMap := make(map[int]*entity.MonsterSourceData)
	for _, monster := range allMonsters {
		monsterMap[monster.No] = monster
	}

	// 対象のモンスターを取得する
	targetMonster, exists := monsterMap[targetNo]
	if !exists {
		return nil, fmt.Errorf("monster not found: %d", targetNo)
	}

	// 起源モンスターまで遡る
	visited := make(map[int]bool) // 循環参照防止
	for targetMonster != nil {
		if isOriginMonster(targetMonster) {
			break
		}
		if visited[targetMonster.No] {
			// 循環検出時の対応
			return nil, fmt.Errorf("circular reference detected: %d", targetMonster.No)
		}
		visited[targetMonster.No] = true
		nextNo := targetMonster.BaseNo
		if nextNo == nil {
			return nil, fmt.Errorf("baseNo not found: %d", targetMonster.No)
		} else {
			targetMonster = monsterMap[*nextNo]
		}
	}

	if targetMonster == nil {
		return nil, fmt.Errorf("origin monster not found: %d", targetNo)
	}

	return targetMonster, nil
}

func isOriginMonster(targetMonster *entity.MonsterSourceData) bool {
	return targetMonster.BaseNo == nil
}
