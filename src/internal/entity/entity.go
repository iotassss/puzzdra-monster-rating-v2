package entity

import "net/url"

type Game8MonsterScore struct {
	Name        string
	LeaderPoint string
	SubPoint    string
	AssistPoint string
}

type Game8Monster struct {
	No     int
	URL    string
	Scores []Game8MonsterScore
}

type GamewithMonster struct {
}

type Monster struct {
	No              int
	Name            string
	OriginMonsterNo int
	Game8URL        *url.URL
	Game8Scores     []Game8MonsterScore
}

type MonsterSourceData struct {
	No     int
	Name   string
	BaseNo int
}
