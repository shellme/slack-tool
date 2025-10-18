# slack-tool 管理用Makefile

.PHONY: build install clean test help

# ビルド
build:
	@echo "slack-tool をビルドしています..."
	go build -o slack-tool ./cmd/slack-tool
	@echo "ビルド完了！"

# インストール
install:
	@echo "slack-tool をインストールしています..."
	go install ./cmd/slack-tool
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
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build -o dist/slack-tool-darwin-amd64 ./cmd/slack-tool
	GOOS=darwin GOARCH=arm64 go build -o dist/slack-tool-darwin-arm64 ./cmd/slack-tool
	GOOS=linux GOARCH=amd64 go build -o dist/slack-tool-linux-amd64 ./cmd/slack-tool
	GOOS=windows GOARCH=amd64 go build -o dist/slack-tool-windows-amd64.exe ./cmd/slack-tool
	@echo "リリース用ビルド完了！"

# ヘルプ
help:
	@echo "slack-tool の管理コマンド:"
	@echo "  make build     - slack-toolをビルド"
	@echo "  make install   - slack-toolをインストール"
	@echo "  make test      - テストを実行"
	@echo "  make clean     - ビルドファイルを削除"
	@echo "  make release   - リリース用ビルド（複数プラットフォーム）"
