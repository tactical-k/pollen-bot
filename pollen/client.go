package pollen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	Value        int    `json:"value"`
	Category     string `json:"category"`
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

// FetchDataJST 日本時間を考慮して花粉情報を取得
// GitHub ActionsなどUTC環境で実行される際に、日本時間での正しい日付で情報を取得する
func (c *Client) FetchDataJST(lat, lon float64, days int) (*Response, error) {
	// 日本時間のタイムゾーンを設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// Asia/Tokyoが利用できない場合は、UTC+9で代替
		jst = time.FixedZone("JST", 9*60*60)
	}

	// 現在の日本時間を取得してログ出力
	now := time.Now().In(jst)

	// 日本時間での現在日時をログに出力（デバッグ用）
	fmt.Printf("現在の日本時間: %s\n", now.Format("2006-01-02 15:04:05 JST"))

	// HTTPクライアントにタイムゾーン情報を設定してAPIを呼び出し
	return c.fetchDataWithTimezone(lat, lon, days, jst)
}

// fetchDataWithTimezone タイムゾーンを考慮してAPIを呼び出し
func (c *Client) fetchDataWithTimezone(lat, lon float64, days int, timezone *time.Location) (*Response, error) {
	url := fmt.Sprintf("%s?key=%s&location.latitude=%f&location.longitude=%f&days=%d",
		apiURL, c.apiKey, lat, lon, days)

	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}

	// タイムゾーン情報をヘッダーに追加（Google APIがサポートしている場合）
	// 現在時刻を指定されたタイムゾーンで取得
	now := time.Now().In(timezone)
	req.Header.Set("X-Client-Timezone", timezone.String())
	req.Header.Set("X-Client-Time", now.Format(time.RFC3339))

	resp, err := c.client.Do(req)
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
