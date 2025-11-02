package main

import (
	"log"
	"os"

	"pollen-discord-bot/notification"
	"pollen-discord-bot/pollen"
	"pollen-discord-bot/util"

	"github.com/joho/godotenv"
)

func main() {
	// .envファイルから環境変数を読み込む（存在する場合のみ）
	// GitHub Actionsなどの本番環境では環境変数が直接設定されるため、
	// .envファイルが存在しなくてもエラーにしない
	_ = godotenv.Load()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if apiKey == "" || webhookURL == "" {
		log.Fatal("環境変数 GOOGLE_API_KEY と DISCORD_WEBHOOK_URL を設定してください")
	}

	// 花粉情報クライアントを作成
	pollenClient := pollen.NewClient(apiKey)

	// 花粉情報を取得（東京の1日分）
	pollenData, err := pollenClient.FetchData(pollen.TokyoLat, pollen.TokyoLon, 1)
	if err != nil {
		log.Fatalf("花粉情報の取得に失敗しました: %v", err)
	}

	// Discord通知クライアントを作成
	discordNotifier := notification.NewDiscordNotifier(webhookURL)

	// Discord通知を送信
	if err := discordNotifier.SendPollenInfo(pollenData, "東京", util.FormatDate); err != nil {
		log.Fatalf("Discord通知の送信に失敗しました: %v", err)
	}

	log.Println("花粉情報を正常に通知しました")
}
