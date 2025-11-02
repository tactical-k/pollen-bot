# 東京花粉情報Discord通知Bot

毎朝7時に東京の花粉情報をDiscordに自動通知するbotです。

## 機能

- Google Pollen APIから東京の花粉情報を取得
- レベルに応じて色分けされたDiscord通知
- GitHub Actionsで完全自動実行（サーバー不要）

## セットアップ手順

### 1. Google Pollen APIキーの取得

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. プロジェクトを作成（または既存のプロジェクトを選択）
3. **API とサービス** → **ライブラリ** から「Pollen API」を検索して有効化
4. **API とサービス** → **認証情報** → **認証情報を作成** → **APIキー** を選択
5. 作成されたAPIキーをコピー

### 2. Discord Webhook URLの取得

1. 通知を送りたいDiscordチャンネルを開く
2. チャンネル設定（歯車アイコン）→ **連携サービス** → **ウェブフック**
3. **新しいウェブフック** を作成
4. ウェブフックURLをコピー

### 3. GitHubリポジトリの設定

1. このコードをGitHubリポジトリにプッシュ
2. リポジトリの **Settings** → **Secrets and variables** → **Actions**
3. 以下の2つのシークレットを追加：
   - `GOOGLE_API_KEY`: Google Pollen APIキー
   - `DISCORD_WEBHOOK_URL`: Discord Webhook URL

### 4. 通知時刻の変更（オプション）

`.github/workflows/pollen-notification.yml` のcron設定を編集：

```yaml
schedule:
  - cron: '0 22 * * *'  # 毎朝7時(JST) = 22時(UTC前日)
```

時刻の変換例：
- 6時(JST) = `'0 21 * * *'` (UTC)
- 7時(JST) = `'0 22 * * *'` (UTC)
- 8時(JST) = `'0 23 * * *'` (UTC)

## 手動実行

GitHub Actionsの **Actions** タブから「Daily Pollen Notification」を選択し、
**Run workflow** で手動実行できます。

## ディレクトリ構成

```
.
├── .github/
│   └── workflows/
│       └── pollen-notification.yml  # GitHub Actionsワークフロー
├── pollen/
│   ├── client.go                    # Pollen APIクライアント
│   └── client_test.go               # APIクライアントのテスト
├── notification/
│   ├── discord.go                   # Discord通知機能
│   └── discord_test.go              # Discord通知のテスト
├── util/
│   ├── formatter.go                 # ユーティリティ関数
│   └── formatter_test.go            # ユーティリティのテスト
├── main.go                          # メインプログラム
├── go.mod                           # Goモジュール定義
└── README.md                        # このファイル
```

## テスト

### テストの実行

すべてのテストを実行：
```bash
go test ./...
```

詳細な出力で実行：
```bash
go test ./... -v
```

カバレッジを確認：
```bash
go test ./... -cover
```

特定のパッケージのみテスト：
```bash
go test ./notification -v
go test ./pollen -v
go test ./util -v
```

### テスト構成

- **pollen/client_test.go** - Pollen APIクライアントのテスト
  - クライアント初期化
  - 座標定数の妥当性チェック

- **notification/discord_test.go** - Discord通知機能のテスト
  - 正常な通知送信
  - シーズン外花粉の処理
  - エラーハンドリング
  - 絵文字と色の選択ロジック

- **util/formatter_test.go** - ユーティリティ関数のテスト
  - 日付フォーマット処理
  - エッジケースの処理

## 花粉レベルの見方

- ✅/🟢 レベル0-2: 少ない
- 🟡 レベル3: やや多い
- 🟠 レベル4: 多い
- 🔴 レベル5: 非常に多い

## 注意事項

- Google Pollen APIは月1,000リクエストまで無料
- 1日1回の実行なら月30リクエスト程度なので無料枠内で運用可能
- GitHub Actionsも無料枠内で十分利用可能

## トラブルシューティング

### 通知が来ない場合

1. GitHub Actionsの **Actions** タブでワークフロー実行ログを確認
2. Secretsが正しく設定されているか確認
3. 手動実行（Run workflow）でテスト

### APIエラーが出る場合

- Google Pollen APIが有効化されているか確認
- APIキーが正しいか確認
- APIの利用制限に達していないか確認

## ライセンス

MIT
