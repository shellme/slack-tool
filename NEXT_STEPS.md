# slack-tool 移行完了 - 次のステップ

## 移行完了状況

✅ ファイル構造の準備完了
✅ 設定ディレクトリパスの変更完了 (`~/.config/chore-bot/` → `~/.config/slack-tool/`)
✅ importパスの更新完了 (`github.com/shellme/chore-bot` → `github.com/shellme/slack-tool`)
✅ GitHub Actionsワークフローの配置完了
✅ 依存関係の整理完了
✅ ビルドとテストの確認完了

## 次に必要な作業

### 1. GitHub上でのリポジトリ作成

#### 1.1 メインリポジトリの作成
```bash
# GitHub上で新しいリポジトリを作成
# リポジトリ名: slack-tool
# 説明: "Slackの様々な操作を行うCLIツール"
# ライセンス: MIT
# README: チェックを外す（既に準備済み）
```

#### 1.2 Homebrew tapリポジトリの作成
```bash
# GitHub上で新しいリポジトリを作成
# リポジトリ名: homebrew-slack-tool
# 説明: "Homebrew tap for slack-tool"
# プライベート: いいえ
```

### 2. ローカルリポジトリの初期化とプッシュ

```bash
cd /Users/hosogaimiki/dev-private/slack-tool

# Gitリポジトリの初期化
git init
git add .
git commit -m "Initial commit: migrate slack-tool from chore-bot"

# リモートリポジトリの追加（GitHub上で作成後）
git remote add origin https://github.com/shellme/slack-tool.git
git branch -M main
git push -u origin main
```

### 3. 初回リリース

```bash
# 初回リリースタグの作成
git tag v0.1.1
git push origin v0.1.1
```

### 4. GitHub Actionsの確認

1. GitHub Actionsが自動的にリリースをビルド
2. [Releases](https://github.com/shellme/slack-tool/releases) ページで確認
3. バイナリとtarballが生成されることを確認

### 5. Homebrew Formulaの更新

#### 5.1 tarballのSHA256を計算
```bash
# リリース後にtarballのSHA256を計算
curl -L https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz | shasum -a 256
```

#### 5.2 Formulaファイルの更新
`Formula/slack-tool.rb`のSHA256を更新:
```ruby
url "https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz"
sha256 "計算されたSHA256値"  # ここを更新
```

#### 5.3 homebrew-slack-toolリポジトリの準備
```bash
# homebrew-slack-toolリポジトリをクローン
git clone https://github.com/shellme/homebrew-slack-tool.git
cd homebrew-slack-tool

# Formulaファイルを配置
cp /Users/hosogaimiki/dev-private/slack-tool/Formula/slack-tool.rb ./Formula/slack-tool.rb
cp /Users/hosogaimiki/dev-private/slack-tool/homebrew-tap-README.md ./README.md

# コミット・プッシュ
git add .
git commit -m "Add slack-tool formula"
git push origin main
```

### 6. インストールテスト

```bash
# Homebrewでインストールテスト
brew tap shellme/slack-tool
brew install slack-tool
slack-tool --version

# 期待される出力: slack-tool version v0.1.1
```

### 7. 既存ユーザーへの通知

`chore-bot`リポジトリのREADMEを更新して移行を通知:

```markdown
# 🚀 slack-tool が独立リポジトリに移行しました！

slack-tool は `github.com/shellme/slack-tool` に独立したリポジトリとして移行しました。

## 新しいインストール方法

### Homebrew（推奨）
```bash
brew tap shellme/slack-tool
brew install slack-tool
```

### Go install
```bash
go install github.com/shellme/slack-tool@latest
```

## 移行のメリット

- 🏠 Homebrewでの簡単インストール
- 🔄 自動更新のサポート
- 📦 独立したバージョン管理
- 🚀 より高速なリリースサイクル

詳細は [slack-tool リポジトリ](https://github.com/shellme/slack-tool) をご確認ください。
```

## 移行後の確認事項

- [ ] GitHub上でリポジトリが作成されている
- [ ] 初回リリースが成功している
- [ ] Homebrewでインストールできる
- [ ] 設定ファイルが正しいパスに保存される
- [ ] 既存ユーザーに移行を通知済み

## トラブルシューティング

### よくある問題

1. **ビルドエラー**
   ```bash
   # 依存関係を再整理
   go mod tidy
   go mod download
   ```

2. **Homebrewインストールの失敗**
   ```bash
   # FormulaファイルのSHA256が正しいか確認
   brew audit --strict slack-tool
   ```

3. **設定ファイルのパス**
   ```bash
   # 設定ディレクトリが正しく変更されているか確認
   grep -r "slack-tool" internal/config/
   ```
