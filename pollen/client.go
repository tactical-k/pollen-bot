package pollen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL = "https://pollen.googleapis.com/v1/forecast:lookup"
	// 東京の座標（東京都庁付近）
	TokyoLat = 35.6762
	TokyoLon = 139.6503
)

// Response Google Pollen APIのレスポンス構造
type Response struct {
	DailyInfo []DailyInfo `json:"dailyInfo"`
}

type DailyInfo struct {
	Date        Date         `json:"date"`
	PollenTypes []PollenType `json:"pollenTypeInfo"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// String 日付を文字列形式（YYYY-MM-DD）で返す
func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

type PollenType struct {
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
	IndexInfo   Index  `json:"indexInfo"`
	InSeason    bool   `json:"inSeason"`
}

type Index struct {
	Value       int    `json:"value"`
	Category    string `json:"category"`
	IndexDisplay string `json:"indexDisplay"`
}

// Client Pollen APIクライアント
type Client struct {
	apiKey string
	client *http.Client
}

// NewClient 新しいPollen APIクライアントを作成
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// FetchData 指定した緯度経度の花粉情報を取得
func (c *Client) FetchData(lat, lon float64, days int) (*Response, error) {
	url := fmt.Sprintf("%s?key=%s&location.latitude=%f&location.longitude=%f&days=%d",
		apiURL, c.apiKey, lat, lon, days)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー (ステータス: %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}

	var pollenResp Response
	if err := json.Unmarshal(body, &pollenResp); err != nil {
		return nil, fmt.Errorf("JSONパース失敗: %w", err)
	}

	return &pollenResp, nil
}
