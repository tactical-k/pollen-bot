package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pollen-discord-bot/pollen"
)

// Webhook Discord Webhookのペイロード
type Webhook struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordNotifier Discord通知を送信するクライアント
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier 新しいDiscordNotifierを作成
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// SendPollenInfo 花粉情報をDiscordに送信
func (d *DiscordNotifier) SendPollenInfo(data *pollen.Response, location string, formatDate func(string) string) error {
	if len(data.DailyInfo) == 0 {
		return fmt.Errorf("花粉データがありません")
	}

	today := data.DailyInfo[0]

	// Embedフィールドを作成
	var fields []Field
	var hasInSeasonPollen bool

	for _, p := range today.PollenTypes {
		if !p.InSeason {
			continue
		}
		hasInSeasonPollen = true

		// レベルに応じた絵文字
		emoji := getPollenEmoji(p.IndexInfo.Value)

		fields = append(fields, Field{
			Name:   fmt.Sprintf("%s %s", emoji, p.DisplayName),
			Value:  fmt.Sprintf("レベル: **%s** (%d)", p.IndexInfo.Category, p.IndexInfo.Value),
			Inline: true,
		})
	}

	// シーズン中の花粉がない場合
	if !hasInSeasonPollen {
		fields = append(fields, Field{
			Name:   "🌸 花粉情報",
			Value:  "現在、シーズン中の花粉はありません",
			Inline: false,
		})
	}

	// 色を決定（最大レベルに基づく）
	color := getEmbedColor(today.PollenTypes)

	webhook := Webhook{
		Embeds: []Embed{
			{
				Title:       fmt.Sprintf("🌸 %sの花粉情報", location),
				Description: fmt.Sprintf("📅 %s の花粉情報", formatDate(today.Date.String())),
				Color:       color,
				Fields:      fields,
			},
		},
	}

	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("JSON作成失敗: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Webhook送信失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Webhookエラー (ステータス: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func getPollenEmoji(level int) string {
	switch {
	case level == 0:
		return "✅"
	case level <= 2:
		return "🟢"
	case level <= 3:
		return "🟡"
	case level <= 4:
		return "🟠"
	default:
		return "🔴"
	}
}

func getEmbedColor(pollens []pollen.PollenType) int {
	maxLevel := 0
	for _, p := range pollens {
		if p.InSeason && p.IndexInfo.Value > maxLevel {
			maxLevel = p.IndexInfo.Value
		}
	}

	// Discord色コード（16進数を10進数に変換）
	switch {
	case maxLevel == 0:
		return 0x00FF00 // 緑
	case maxLevel <= 2:
		return 0x00FF00 // 緑
	case maxLevel <= 3:
		return 0xFFFF00 // 黄色
	case maxLevel <= 4:
		return 0xFFA500 // オレンジ
	default:
		return 0xFF0000 // 赤
	}
}
