package service

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
var timeoutSecond = 5

type Game8MonsterScoreScrapingResult struct {
	Name        string
	LeaderPoint string
	SubPoint    string
	AssistPoint string
}

type Game8MonsterScrapingResult struct {
	// Noには何かしらのそのページのモンスターに紐づくNoが入るはず
	No     int
	URL    string
	Scores []*Game8MonsterScoreScrapingResult
}

// type gamewithMonsterScrapingResult struct {
// }
