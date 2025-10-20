# 開発者向け情報

> [!NOTE]
> このドキュメントは、slack-toolの開発・ビルド・リリースに関する情報を提供します。一般ユーザーは[メインREADME](../README.md)を参照してください。

## 前提条件

- Go 1.21 以上
- Git
- Make（オプション）

## ビルド

```bash
# リポジトリをクローン
git clone https://github.com/shellme/slack-tool.git
cd slack-tool

# 依存関係をインストール
go mod download

# ビルド
go build -o slack-tool ./cmd/slack-tool

# または Makefile を使用
make build
```

## 開発ワークフロー

### 新機能の追加

> [!IMPORTANT]
> 新機能を追加する際は、必ずテストを記述し、ビルドが成功することを確認してください。

1. 機能ブランチを作成
   ```bash
   git checkout -b feature/new-feature
   ```

2. 機能を実装

3. テストを実行
   ```bash
   go test ./...
   ```

4. ビルドを確認
   ```bash
   make build
   ```

5. コミット・プッシュ
   ```bash
   git add .
   git commit -m "feat: 新機能の説明"
   git push origin feature/new-feature
   ```

6. プルリクエストを作成

### リリース

1. バージョンを更新
   ```bash
   # バージョンタグを作成
   git tag v0.1.2
   git push origin v0.1.2
   ```

2. Homebrew Formulaを更新
   - `Formula/slack-tool.rb` のバージョンとSHAを更新

3. GitHub Releaseを作成

## プロジェクト構造

```
slack-tool/
├── cmd/slack-tool/          # メインアプリケーション
│   ├── cmd/                 # コマンド定義
│   │   ├── channel.go       # チャンネル取得コマンド
│   │   ├── config.go        # 設定コマンド
│   │   ├── get.go           # データ取得コマンド
│   │   ├── post.go          # メッセージ投稿コマンド
│   │   ├── reactions.go     # リアクション取得コマンド
│   │   └── root.go          # ルートコマンド
│   └── main.go              # エントリーポイント
├── internal/                # 内部パッケージ
│   ├── config/              # 設定管理
│   └── slack/               # Slack API クライアント
├── docs/                    # ドキュメント
├── Formula/                 # Homebrew Formula
├── go.mod                   # Go モジュール定義
├── go.sum                   # 依存関係のチェックサム
├── Makefile                 # ビルドスクリプト
└── README.md                # メインドキュメント
```

## コーディング規約

- Go標準のフォーマットを使用
- コメントは日本語で記述
- エラーハンドリングは適切に行う
- テストは可能な限り記述

## 依存関係

- [cobra](https://github.com/spf13/cobra) - CLI フレームワーク
- [slack-go/slack](https://github.com/slack-go/slack) - Slack API クライアント

## ライセンス

MIT License
