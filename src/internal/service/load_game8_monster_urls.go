package service

import (
	"bufio"
	"net/url"
	"os"
)

func LoadGame8MonsterURLs(game8MonsterURLListFilePath string) ([]string, error) {
	// ファイルを開く
	file, err := os.Open(game8MonsterURLListFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string

	// ファイルを1行ずつ読み込む
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// URLが有効か確認
		_, err := url.ParseRequestURI(line)
		if err != nil {
			return nil, err // 無効なURLがあればエラーを返す
		}

		// 有効なURLをリストに追加
		urls = append(urls, line)
	}

	// スキャン中にエラーが発生した場合
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
