# slack-tool 管理用Makefile

.PHONY: build install clean test help

# バージョン情報の取得
# オプション1: 完全な情報（現在の設定）
VERSION := $(shell git describe --tags --always --dirty)
# オプション2: シンプルなタグのみ
# VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
# オプション3: タグ + コミット数（dirtyフラグなし）
# VERSION := $(shell git describe --tags --always)

COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# ビルドフラグ
LDFLAGS := -ldflags "-X github.com/shellme/slack-tool/cmd/slack-tool/cmd.version=$(VERSION) -X github.com/shellme/slack-tool/cmd/slack-tool/cmd.commit=$(COMMIT) -X github.com/shellme/slack-tool/cmd/slack-tool/cmd.date=$(DATE)"

# ビルド
build:
	@echo "slack-tool をビルドしています..."
	@echo "バージョン: $(VERSION)"
	@echo "コミット: $(COMMIT)"
	@echo "日時: $(DATE)"
	go build $(LDFLAGS) -o slack-tool ./cmd/slack-tool
	@echo "ビルド完了！"

# インストール
install:
	@echo "slack-tool をインストールしています..."
	@echo "バージョン: $(VERSION)"
	@echo "コミット: $(COMMIT)"
	@echo "日時: $(DATE)"
	go install $(LDFLAGS) ./cmd/slack-tool
	@echo "インストール完了！"

# テスト
test:
	@echo "テストを実行しています..."
	go test ./...

# クリーンアップ
clean:
	@echo "ビルドファイルを削除しています..."
	rm -f slack-tool
	@echo "クリーンアップ完了！"

# リリース用ビルド（複数プラットフォーム）
release:
	@echo "リリース用ビルドを実行しています..."
	@echo "バージョン: $(VERSION)"
	@echo "コミット: $(COMMIT)"
	@echo "日時: $(DATE)"
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/slack-tool-darwin-amd64 ./cmd/slack-tool
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/slack-tool-darwin-arm64 ./cmd/slack-tool
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/slack-tool-linux-amd64 ./cmd/slack-tool
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/slack-tool-windows-amd64.exe ./cmd/slack-tool
	@echo "リリース用ビルド完了！"

# ヘルプ
help:
	@echo "slack-tool の管理コマンド:"
	@echo "  make build     - slack-toolをビルド"
	@echo "  make install   - slack-toolをインストール"
	@echo "  make test      - テストを実行"
	@echo "  make clean     - ビルドファイルを削除"
	@echo "  make release   - リリース用ビルド（複数プラットフォーム）"
