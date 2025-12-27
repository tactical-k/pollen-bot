package notification

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"pollen-discord-bot/pollen"
)

func TestNewDiscordNotifier(t *testing.T) {
	webhookURL := "https://discord.com/api/webhooks/test"
	notifier := NewDiscordNotifier(webhookURL)

	if notifier.webhookURL != webhookURL {
		t.Errorf("expected webhookURL %s, got %s", webhookURL, notifier.webhookURL)
	}

	if notifier.client == nil {
		t.Error("expected http.Client to be initialized")
	}
}

func TestSendPollenInfo_Success(t *testing.T) {
	// モックサーバーを作成
	var receivedPayload Webhook
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// テストデータを作成
	testData := &pollen.Response{
		DailyInfo: []pollen.DailyInfo{
			{
				Date: pollen.Date{Year: 2024, Month: 3, Day: 15},
				PollenTypes: []pollen.PollenType{
					{
						Code:        "TREE",
						DisplayName: "Tree",
						IndexInfo: pollen.Index{
							Value:    3,
							Category: "MODERATE",
						},
						InSeason: true,
					},
					{
						Code:        "GRASS",
						DisplayName: "Grass",
						IndexInfo: pollen.Index{
							Value:    2,
							Category: "LOW",
						},
						InSeason: true,
					},
				},
			},
		},
	}

	// テスト用のフォーマット関数
	formatDate := func(date string) string {
		return date + "_formatted"
	}

	notifier := NewDiscordNotifier(server.URL)
	err := notifier.SendPollenInfo(testData, "Test Location", formatDate)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// ペイロードの検証
	if len(receivedPayload.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(receivedPayload.Embeds))
	}

	embed := receivedPayload.Embeds[0]
	if embed.Title != "🌸 Test Locationの花粉情報" {
		t.Errorf("unexpected title: %s", embed.Title)
	}

	if len(embed.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(embed.Fields))
	}
}

func TestSendPollenInfo_NoInSeasonPollen(t *testing.T) {
	// モックサーバーを作成
	var receivedPayload Webhook
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// シーズン外のデータ
	testData := &pollen.Response{
		DailyInfo: []pollen.DailyInfo{
			{
				Date: pollen.Date{Year: 2024, Month: 3, Day: 15},
				PollenTypes: []pollen.PollenType{
					{
						Code:        "TREE",
						DisplayName: "Tree",
						IndexInfo: pollen.Index{
							Value:    0,
							Category: "NONE",
						},
						InSeason: false,
					},
				},
			},
		},
	}

	formatDate := func(date string) string { return date }
	notifier := NewDiscordNotifier(server.URL)
	err := notifier.SendPollenInfo(testData, "Test Location", formatDate)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// シーズン外メッセージの確認
	embed := receivedPayload.Embeds[0]
	if len(embed.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(embed.Fields))
	}

	if embed.Fields[0].Value != "現在、シーズン中の花粉はありません" {
		t.Errorf("unexpected message: %s", embed.Fields[0].Value)
	}
}

func TestSendPollenInfo_EmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	testData := &pollen.Response{
		DailyInfo: []pollen.DailyInfo{},
	}

	formatDate := func(date string) string { return date }
	notifier := NewDiscordNotifier(server.URL)
	err := notifier.SendPollenInfo(testData, "Test Location", formatDate)

	if err == nil {
		t.Error("expected error for empty data, got nil")
	}
}

func TestSendPollenInfo_WebhookError(t *testing.T) {
	// エラーを返すモックサーバー
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	testData := &pollen.Response{
		DailyInfo: []pollen.DailyInfo{
			{
				Date:        pollen.Date{Year: 2024, Month: 3, Day: 15},
				PollenTypes: []pollen.PollenType{},
			},
		},
	}

	formatDate := func(date string) string { return date }
	notifier := NewDiscordNotifier(server.URL)
	err := notifier.SendPollenInfo(testData, "Test Location", formatDate)

	if err == nil {
		t.Error("expected error for webhook failure, got nil")
	}
}

func TestGetPollenEmoji(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{0, "✅"},
		{1, "🟢"},
		{2, "🟢"},
		{3, "🟡"},
		{4, "🟠"},
		{5, "🔴"},
	}

	for _, tt := range tests {
		result := getPollenEmoji(tt.level)
		if result != tt.expected {
			t.Errorf("getPollenEmoji(%d) = %s, want %s", tt.level, result, tt.expected)
		}
	}
}

func TestGetEmbedColor(t *testing.T) {
	tests := []struct {
		name     string
		pollens  []pollen.PollenType
		expected int
	}{
		{
			name: "No in-season pollen",
			pollens: []pollen.PollenType{
				{InSeason: false, IndexInfo: pollen.Index{Value: 5}},
			},
			expected: 0x00FF00,
		},
		{
			name: "Low level",
			pollens: []pollen.PollenType{
				{InSeason: true, IndexInfo: pollen.Index{Value: 2}},
			},
			expected: 0x00FF00,
		},
		{
			name: "Moderate level",
			pollens: []pollen.PollenType{
				{InSeason: true, IndexInfo: pollen.Index{Value: 3}},
			},
			expected: 0xFFFF00,
		},
		{
			name: "High level",
			pollens: []pollen.PollenType{
				{InSeason: true, IndexInfo: pollen.Index{Value: 4}},
			},
			expected: 0xFFA500,
		},
		{
			name: "Very high level",
			pollens: []pollen.PollenType{
				{InSeason: true, IndexInfo: pollen.Index{Value: 5}},
			},
			expected: 0xFF0000,
		},
		{
			name: "Multiple pollens - max level",
			pollens: []pollen.PollenType{
				{InSeason: true, IndexInfo: pollen.Index{Value: 2}},
				{InSeason: true, IndexInfo: pollen.Index{Value: 5}},
				{InSeason: true, IndexInfo: pollen.Index{Value: 3}},
			},
			expected: 0xFF0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getEmbedColor(tt.pollens)
			if result != tt.expected {
				t.Errorf("getEmbedColor() = 0x%X, want 0x%X", result, tt.expected)
			}
		})
	}
}
