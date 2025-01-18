package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/stacktrace"
)

func ScrapeGame8(url *url.URL) (Game8MonsterScrapingResult, error) {
	result := Game8MonsterScrapingResult{
		URL: url.String(),
	}

	c := colly.NewCollector(
		colly.UserAgent(userAgent),
	)
	c.SetRequestTimeout(time.Duration(timeoutSecond) * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Scraping error", slog.Any("error", err))
	})

	// No取得
	c.OnHTML("h3", func(e *colly.HTMLElement) {
		// h3要素が「のステータス」を含むことを確認
		if !strings.Contains(e.Text, "のステータス") {
			return
		}

		// h3タグの次のsiblingがtableかどうかを確認
		table := e.DOM.Next()
		if goquery.NodeName(table) != "table" {
			return
		}

		// tableの最初の行のth要素のテキストを確認
		thText := table.Find("tr:first-child th").Text()
		if thText == "" {
			return
		}

		// 【No.xxx】モンスター名 の形式からxxx部分の数値を正規表現で抽出
		re := regexp.MustCompile(`【No\.\s*(\d+)】`)
		match := re.FindStringSubmatch(thText)
		if len(match) < 2 {
			return
		}
		noStr := match[1]

		// 保存
		no, err := strconv.Atoi(noStr)
		if err != nil {
			slog.Error("Failed to convert string to int", slog.String("noStr", noStr), slog.Any("error", err))
			return
		}

		result.No = no
	})

	// 点数取得
	c.OnHTML("table", func(e *colly.HTMLElement) {
		var isScoreTable bool
		var name string
		var isPattern2 bool
		e.ForEach("tr", func(index int, row *colly.HTMLElement) {
			// 最初の行に「リーダー評価」「サブ評価」「アシスト評価」の文字列があるか確認
			isScoreTableHeader := strings.Contains(row.Text, "リーダー") &&
				strings.Contains(row.Text, "サブ") &&
				strings.Contains(row.Text, "アシスト")

			if index == 0 && isScoreTableHeader {
				if strings.Contains(row.Text, "リーダー評価") {
					isPattern2 = true
				}
				isScoreTable = true
				return
			}

			// リーダー評価などが確認されたら次の行から点数を取得
			if isScoreTable {
				var leader, sub, assist string
				if isPattern2 {
					table := row.DOM.Parent().Parent()
					p := table.Prev()
					h2 := p.Prev()

					if goquery.NodeName(h2) != "h2" || !strings.Contains(h2.Text(), "の評価") {
						return
					}

					name = strings.Replace(h2.Text(), "の評価", "", -1)
					name = strings.Replace(name, "と使い道", "", -1)
					leader = row.ChildText("td:nth-of-type(1)")
					leader = strings.TrimSuffix(leader, "点 / 9.9点")
					sub = row.ChildText("td:nth-of-type(2)")
					sub = strings.TrimSuffix(sub, "点 / 9.9点")
					assist = row.ChildText("td:nth-of-type(3)")
					assist = strings.TrimSuffix(assist, "点 / 9.9点")
				} else {
					name = row.ChildText("td:nth-of-type(1)")
					leader = row.ChildText("td:nth-of-type(2)")
					sub = row.ChildText("td:nth-of-type(3)")
					assist = row.ChildText("td:nth-of-type(4)")
				}

				result.Scores = append(result.Scores, &Game8MonsterScoreScrapingResult{
					Name:        name,
					LeaderPoint: leader,
					SubPoint:    sub,
					AssistPoint: assist,
				})
			}
		})
	})

	err := c.Visit(url.String())
	if err != nil {
		slog.Error("Failed to visit url", slog.String("url", url.String()), slog.Any("error", err), slog.String("stacktrace", stacktrace.Print()))
		return Game8MonsterScrapingResult{}, err
	}

	// もしこの時点でresult.noが0ならばNoが取得できていないのでエラー
	if result.No == 0 {
		err := fmt.Errorf("Failed to scrape game8 monster")
		slog.Warn("Failed to scrape game8 monster", slog.String("url", url.String()), slog.Any("error", err), slog.String("stacktrace", stacktrace.Print()))
		return Game8MonsterScrapingResult{}, err
	}

	return result, nil
}
