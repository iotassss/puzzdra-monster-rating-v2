package batch

import (
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// モンスターurl一覧取得処理を実行しgame8_monster_urls.txtを作成
func (batch *Batch) CollectGame8MonsterURLs() error {
	file, err := os.Create(batch.game8MonsterURLListFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	c := colly.NewCollector()

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		// 副属性「*」のパターンを正規表現に変換
		pattern := `^副属性「.*」のキャラ評価$`
		re := regexp.MustCompile(pattern)
		if !re.MatchString(e.Text) {
			return
		}

		// h3タグの次にあるtableを取得
		nextSibling := e.DOM.Next()
		if goquery.NodeName(nextSibling) != "table" {
			return
		}

		// table内のリンクを取得
		nextSibling.Find("tr").Each(func(i int, selection *goquery.Selection) {
			selection.Find("td a").Each(func(j int, a *goquery.Selection) {
				href, exists := a.Attr("href")
				if exists {
					file.WriteString(href + "\n")
				}
			})
		})
	})

	// モンスター一覧ページのURL
	monsterListUrls := []string{
		"https://game8.jp/pazudora/24173", // 火属性
		"https://game8.jp/pazudora/24241", // 水属性
		"https://game8.jp/pazudora/24242", // 木属性
		"https://game8.jp/pazudora/24243", // 光属性
		"https://game8.jp/pazudora/24236", // 闇属性
	}

	for _, url := range monsterListUrls {
		c.Visit(url)
	}

	return nil
}
