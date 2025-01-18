package service_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/iotassss/puzzdra-monster-rating-v2/internal/service"
	"github.com/stretchr/testify/assert"
)

// オメガモン
func TestScrapeGame8(t *testing.T) {
	// HTMLファイルの内容を読み込み
	omegamonHtmlContent, err := os.ReadFile("../../../src/test/testdata/game8/omegamon.html")
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}
	amaterasuHtmlContent, err := os.ReadFile("../../../src/test/testdata/game8/amaterasu.html")
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	// Mockサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/omegamon":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(omegamonHtmlContent)
		case r.URL.Path == "/amaterasu":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(amaterasuHtmlContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	t.Run("Omegamon", func(t *testing.T) {
		// // Mockサーバーへのリクエスト
		// resp, err := http.Get(mockServer.URL + "/omegamon")
		// if err != nil {
		// 	t.Fatalf("Request failed: %v", err)
		// }
		// defer resp.Body.Close()

		// // レスポンスの確認
		// if resp.StatusCode != http.StatusOK {
		// 	t.Errorf("Expected status 200, got %d", resp.StatusCode)
		// }

		// // read the response body
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	t.Fatalf("Failed to read response body: %v", err)
		// }

		// // レスポンスボディの確認
		// if string(body) != string(omegamonHtmlContent) {
		// 	t.Errorf("Expected body %s, got %s", string(omegamonHtmlContent), string(body))
		// }

		testURL, err := url.Parse(mockServer.URL + "/omegamon")
		assert.NoError(t, err)

		result, err := service.ScrapeGame8(testURL)
		assert.NoError(t, err)

		// 保留
		// assert.Equal(t, 11714, result.No)
		assert.Equal(t, testURL.String(), result.URL)
		assert.Equal(t, "10.0", result.Scores[0].LeaderPoint)
		assert.Equal(t, "9.0", result.Scores[0].SubPoint)
		assert.Equal(t, "-", result.Scores[0].AssistPoint)
	})

	t.Run("Amaterasu", func(t *testing.T) {
		testURL, err := url.Parse(mockServer.URL + "/amaterasu")
		assert.NoError(t, err)

		result, err := service.ScrapeGame8(testURL)
		assert.NoError(t, err)

		// 保留
		// assert.Equal(t, 11540, result.No)
		assert.Equal(t, testURL.String(), result.URL)
		assert.Equal(t, "超転生アマテラス", result.Scores[0].Name)
		assert.Equal(t, "5.0", result.Scores[0].LeaderPoint)
		assert.Equal(t, "8.0", result.Scores[0].SubPoint)
		assert.Equal(t, "-", result.Scores[0].AssistPoint)
		assert.Equal(t, "アマテラス装備（櫛）", result.Scores[1].Name)
		assert.Equal(t, "-", result.Scores[1].LeaderPoint)
		assert.Equal(t, "-", result.Scores[1].SubPoint)
		assert.Equal(t, "8.0", result.Scores[1].AssistPoint)
		assert.Equal(t, "アマテラス装備（懐中時計）", result.Scores[2].Name)
		assert.Equal(t, "-", result.Scores[2].LeaderPoint)
		assert.Equal(t, "-", result.Scores[2].SubPoint)
		assert.Equal(t, "7.5", result.Scores[2].AssistPoint)
		assert.Equal(t, "試練アマテラス", result.Scores[3].Name)
		assert.Equal(t, "7.5", result.Scores[3].LeaderPoint)
		assert.Equal(t, "9.0", result.Scores[3].SubPoint)
		assert.Equal(t, "-", result.Scores[3].AssistPoint)
	})
}

// func requestMockServer(url string, htmlContent []byte) ()
