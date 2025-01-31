package handler

import (
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/repository"
)

type Handler struct {
	monsterRepo *repository.MonsterRepository
}

func NewHandler(monsterRepo *repository.MonsterRepository) *Handler {
	return &Handler{monsterRepo: monsterRepo}
}
