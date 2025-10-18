# Homebrewテストコマンド

## 概要

Homebrewでのslack-toolインストールをテストするコマンドです。既存のインストールをアンインストールし、最新バージョンをインストールして動作確認を行います。

## 使用方法

```
/test-homebrew
```

## パラメータ

なし

## 例

```
/test-homebrew
```

## 実行方法

このコマンドが実行された場合、以下のスクリプトを実行してください：
```bash
./.cursor/scripts/test-homebrew.sh
```

## 手動実行（人間向け）

```bash
# 直接スクリプトを実行
./.cursor/scripts/test-homebrew.sh

# ヘルプを表示
./.cursor/scripts/test-homebrew.sh --help
```

## 実行内容

1. **既存インストールの確認**
   - slack-toolがインストールされているかチェック
   - インストールされている場合はアンインストール

2. **tapの確認・追加**
   - shellme/slack-toolのtapが追加されているかチェック
   - 追加されていない場合は自動追加

3. **キャッシュのクリア**
   - Homebrewのキャッシュをクリア
   - 古いバージョンが残らないようにする

4. **新規インストール**
   - 最新バージョンのslack-toolをインストール
   - インストール成功を確認

5. **動作テスト**
   - バージョン確認
   - ヘルプ表示の確認
   - 基本コマンドの動作確認

## 注意事項

- Homebrewがインストールされている必要があります
- 管理者権限が必要な場合があります
- 既存のslack-toolがアンインストールされます

## 前提条件

- Homebrewがインストールされている
- インターネット接続がある
- homebrew-slack-toolリポジトリが公開されている

## トラブルシューティング

### インストール失敗時

```bash
# デバッグモードで実行
brew install --debug slack-tool

# Formulaの構文チェック
brew audit --strict slack-tool

# キャッシュをクリア
brew cleanup slack-tool
brew uninstall slack-tool
brew tap shellme/slack-tool
brew install slack-tool
```

### よくあるエラー

1. **tapが追加されていない**
   - 手動でtapを追加: `brew tap shellme/slack-tool`

2. **インストール失敗**
   - ネットワーク接続を確認
   - Homebrewを更新: `brew update`

3. **バージョンが一致しない**
   - キャッシュをクリア: `brew cleanup`
   - 再インストール: `brew uninstall slack-tool && brew install slack-tool`

4. **権限エラー**
   - 管理者権限で実行
   - Homebrewの権限を確認

5. **Formulaが見つからない**
   - homebrew-slack-toolリポジトリが存在するか確認
   - tapが正しく追加されているか確認

6. **実行ファイルの権限エラー**
   - エラー例: `permission denied: slack-tool`
   - 原因: Formulaのビルドコマンドで実行権限が正しく設定されていない
   - 対処法: Formulaファイルのビルドコマンドを確認・修正

7. **実行ファイルが破損している**
   - エラー例: `syntax error near unexpected token 'newline'`
   - 原因: ディレクトリ構造の問題で正しい実行ファイルが生成されていない
   - 対処法: ローカルでビルドテストを実行して構造を確認

## 期待される動作

### 正常な場合

```
[SUCCESS] slack-tool v0.1.2 のインストールが完了しました
[SUCCESS] バージョン確認: slack-tool version v0.1.2
[SUCCESS] ヘルプ表示: 正常
[SUCCESS] すべてのテストが成功しました！
```

### エラーの場合

```
[ERROR] インストールに失敗しました
[ERROR] バージョンが一致しません
[ERROR] コマンドが実行できません
```

## ヘルパーコマンド

- `/release` - 新バージョンのリリース
- `/update-homebrew` - Homebrew Formulaの更新

## 実行後の確認

テスト完了後、以下のコマンドで動作を確認できます：

```bash
# バージョン確認
slack-tool --version

# ヘルプ表示
slack-tool --help

# 設定確認
slack-tool config show
```

## 自動化のヒント

このコマンドは以下の場面で特に有用です：

- リリース後の動作確認
- Formula更新後のテスト
- 新しい環境でのセットアップ確認
- CI/CDパイプラインでのテスト
