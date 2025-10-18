# Homebrew更新コマンド

## 概要

Homebrew Formulaを更新するコマンドです。指定したバージョンのtarballのSHA256を計算し、Formulaファイルを更新してhomebrew-slack-toolリポジトリに反映します。

## 使用方法

```
/update-homebrew [version]
```

## パラメータ

- `version`: 更新対象のバージョン（省略時は最新リリース）

## 例

```
/update-homebrew v0.1.2
/update-homebrew
```

## 実行方法

このコマンドが実行された場合、以下のスクリプトを実行してください：
```bash
./.cursor/scripts/update-homebrew.sh [version]
```

## 手動実行（人間向け）

```bash
# 直接スクリプトを実行
./.cursor/scripts/update-homebrew.sh v0.1.2

# 最新リリースを更新
./.cursor/scripts/update-homebrew.sh

# ヘルプを表示
./.cursor/scripts/update-homebrew.sh --help
```

## 実行内容

1. **リリース確認**
   - 指定バージョンのリリースが存在するか確認
   - 省略時は最新リリースを取得

2. **SHA256計算**
   - リリースのtarballをダウンロード
   - SHA256ハッシュを計算

3. **Formulaファイル更新**
   - `Formula/slack-tool.rb`のURLとSHA256を更新
   - バージョン番号を更新

4. **homebrew-slack-toolリポジトリに反映**
   - リポジトリをクローン
   - 更新されたFormulaファイルを配置
   - 変更をコミット・プッシュ

## 注意事項

- GitHub CLI (gh) がインストールされている必要があります
- homebrew-slack-toolリポジトリへのプッシュ権限が必要です
- リリースが存在しない場合はエラーになります

## 前提条件

- GitHub CLI (gh) がインストールされている
- shellmeアカウントで認証されている
- homebrew-slack-toolリポジトリが存在する
- プッシュ権限がある

## トラブルシューティング

### Homebrew更新失敗時

```bash
# 手動でSHA256を計算
curl -L https://github.com/shellme/slack-tool/archive/v0.1.2.tar.gz | shasum -a 256

# Formulaファイルを手動編集
vim Formula/slack-tool.rb

# homebrew-slack-toolリポジトリを手動で更新
cd /tmp/homebrew-slack-tool
git add Formula/slack-tool.rb
git commit -m "Update slack-tool to v0.1.2"
git push origin main
```

### よくあるエラー

1. **リリースが存在しない**
   - 指定したバージョンのリリースがGitHub上に存在するか確認
   - 先に `/release` コマンドでリリースを作成してください

2. **tarballのダウンロード失敗**
   - ネットワーク接続を確認
   - リリースURLが正しいか確認

3. **SHA256の計算失敗**
   - ダウンロードしたファイルが破損していないか確認
   - 手動でSHA256を計算して比較

4. **homebrew-slack-toolのクローン失敗**
   - リポジトリが存在するか確認
   - アクセス権限を確認

5. **プッシュ権限がない**
   - GitHub CLIで認証されているか確認
   - リポジトリへの書き込み権限があるか確認

## ヘルパーコマンド

- `/release` - 新バージョンのリリース
- `/test-homebrew` - Homebrewインストールのテスト

## 実行後の確認

更新完了後、以下のコマンドでインストールテストを実行してください：

```bash
/test-homebrew
```

これにより、更新されたFormulaが正しく動作することを確認できます。
