# slack-tool リポジトリ移行ガイド

このガイドでは、`chore-bot` リポジトリから `slack-tool` を独立したリポジトリに移行する手順を説明します。

## 移行の概要

- **現在**: `github.com/shellme/chore-bot` の一部として `slack-tool` が存在
- **移行後**: `github.com/shellme/slack-tool` として独立したリポジトリ

## 移行手順

### 1. 新しいリポジトリの作成

1. GitHubで新しいリポジトリを作成
   - リポジトリ名: `slack-tool`
   - 説明: "Slackの様々な操作を行うCLIツール"
   - ライセンス: MIT
   - README: チェックを外す（後で追加）

2. Homebrew tap用のリポジトリを作成
   - リポジトリ名: `homebrew-slack-tool`
   - 説明: "Homebrew tap for slack-tool"
   - プライベート: いいえ

### 2. ファイルの移行

#### 2.1 slack-toolリポジトリに移行するファイル

```bash
# 新しいリポジトリをクローン
git clone https://github.com/shellme/slack-tool.git
cd slack-tool

# 以下のファイルをコピー
cp -r /path/to/chore-bot/cmd/slack-tool ./
cp -r /path/to/chore-bot/internal ./
cp /path/to/chore-bot/slack-tool-go.mod ./go.mod
cp /path/to/chore-bot/slack-tool-Makefile ./Makefile
cp /path/to/chore-bot/slack-tool-README.md ./README.md
cp /path/to/chore-bot/.github-workflows-release.yml ./.github/workflows/release.yml
```

#### 2.2 設定ファイルの調整

1. **go.modの調整**
   ```bash
   # 依存関係を更新
   go mod tidy
   ```

2. **importパスの更新**
   ```bash
   # 全てのGoファイルでimportパスを更新
   find . -name "*.go" -exec sed -i 's|github.com/shellme/chore-bot|github.com/shellme/slack-tool|g' {} \;
   ```

3. **設定ディレクトリの変更**
   ```bash
   # internal/config/config.go で設定ディレクトリを変更
   # ~/.config/chore-bot/ → ~/.config/slack-tool/
   ```

### 3. Homebrew tapの設定

#### 3.1 homebrew-slack-toolリポジトリの設定

```bash
# homebrew-slack-toolリポジトリをクローン
git clone https://github.com/shellme/homebrew-slack-tool.git
cd homebrew-slack-tool

# Formulaファイルを配置
cp /path/to/chore-bot/slack-tool.rb ./Formula/slack-tool.rb
cp /path/to/chore-bot/homebrew-tap-README.md ./README.md
```

#### 3.2 Formulaファイルの更新

1. **SHA256の計算**
   ```bash
   # リリースタグのtarballのSHA256を計算
   curl -L https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz | shasum -a 256
   ```

2. **Formulaファイルの更新**
   ```ruby
   # slack-tool.rb の url と sha256 を更新
   url "https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz"
   sha256 "計算されたSHA256値"
   ```

### 4. 初回リリース

#### 4.1 slack-toolリポジトリでのリリース

```bash
cd slack-tool

# 初回コミット
git add .
git commit -m "Initial commit: migrate slack-tool from chore-bot"

# 初回リリースタグ
git tag v0.1.1
git push origin main
git push origin v0.1.1
```

#### 4.2 GitHub Actionsの確認

1. GitHub Actionsが自動的にリリースをビルド
2. [Releases](https://github.com/shellme/slack-tool/releases) ページで確認

#### 4.3 Homebrew Formulaの更新

```bash
cd homebrew-slack-tool

# Formulaファイルを更新（SHA256を最新の値に）
# コミット・プッシュ
git add .
git commit -m "Add slack-tool formula"
git push origin main
```

### 5. 既存ユーザーへの通知

#### 5.1 chore-botリポジトリの更新

1. **READMEの更新**
   - 新しいインストール方法を追加
   - 移行に関する注意事項を記載

2. **非推奨警告の追加**
   - 既存のインストール方法に非推奨警告を追加

#### 5.2 移行通知の作成

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

## 移行後のメンテナンス

### 1. バージョン管理

- `slack-tool` リポジトリで独立したバージョン管理
- セマンティックバージョニングに従う
- リリースノートの作成

### 2. Homebrew Formulaの更新

- 新しいリリース時にFormulaファイルを更新
- SHA256の再計算
- テストの実行

### 3. ドキュメントの同期

- 両リポジトリのドキュメントを同期
- 移行完了後は `slack-tool` リポジトリをメインに

## トラブルシューティング

### よくある問題

1. **importパスのエラー**
   ```bash
   # 全てのファイルでimportパスが正しく更新されているか確認
   grep -r "github.com/shellme/chore-bot" .
   ```

2. **設定ファイルのパス**
   ```bash
   # 設定ディレクトリが正しく変更されているか確認
   grep -r "chore-bot" internal/config/
   ```

3. **Homebrewインストールの失敗**
   ```bash
   # FormulaファイルのSHA256が正しいか確認
   brew audit --strict slack-tool
   ```

## 完了チェックリスト

- [ ] 新しいリポジトリの作成
- [ ] ファイルの移行
- [ ] importパスの更新
- [ ] 設定ディレクトリの変更
- [ ] Homebrew Formulaの作成
- [ ] 初回リリース
- [ ] 既存ユーザーへの通知
- [ ] ドキュメントの更新
- [ ] テストの実行
- [ ] 移行の完了確認
