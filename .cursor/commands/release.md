# リリースコマンド

## 概要

slack-toolの新バージョンをリリースするコマンドです。指定したバージョンでGitタグを作成し、GitHub Actionsを実行してリリースを自動化します。

## 使用方法

```
/release [version] [message]
```

## パラメータ

- `version`: リリースバージョン（例: v0.1.2, v0.2.0）
- `message`: リリースメッセージ（省略可、デフォルト: "Release [version]"）

## 例

```
/release v0.1.2 "バグ修正とパフォーマンス改善"
/release v0.2.0
```

## 実行方法

このコマンドが実行された場合、以下のスクリプトを実行してください：
```bash
./.cursor/scripts/release.sh [version] [message]
```

## 手動実行（人間向け）

```bash
# 直接スクリプトを実行
./.cursor/scripts/release.sh v0.1.2 "バグ修正とパフォーマンス改善"

# ヘルプを表示
./.cursor/scripts/release.sh --help
```

## 実行内容

1. **バージョン番号の検証**
   - セマンティックバージョニング形式の確認
   - 既存タグとの重複チェック

2. **環境の確認**
   - 現在のブランチがmainかチェック
   - 未コミット変更の有無を確認

3. **Gitタグの作成**
   - 指定バージョンでタグを作成
   - リリースメッセージをタグに付与

4. **GitHubへのプッシュ**
   - タグをリモートリポジトリにプッシュ
   - GitHub Actionsの実行をトリガー

5. **リリース監視**
   - GitHub Actionsの実行状態を監視
   - リリース成功の確認
   - GitHub Actionsが失敗した場合は手動でリリース作成

## 注意事項

- mainブランチで実行してください
- コミット前にテストを実行してください
- バージョン番号はセマンティックバージョニングに従ってください
- 既存のタグと同じバージョンは使用できません

## 前提条件

- GitHub CLI (gh) がインストールされている
- shellmeアカウントで認証されている
- リモートリポジトリが設定されている
- プッシュ権限がある

## トラブルシューティング

### リリース失敗時

```bash
# タグを削除して再実行
git tag -d v0.1.2
git push origin :refs/tags/v0.1.2

# GitHub Actions のログを確認
gh run list --repo shellme/slack-tool
gh run view <run-id> --log

# 手動でリリースを作成（GitHub Actionsが失敗した場合）
gh release create v0.1.2 --title "v0.1.2" --notes "リリースノート"

# リリースの確認
gh release view v0.1.2
```

### よくあるエラー

1. **バージョン番号が不正**
   - セマンティックバージョニング形式を使用してください
   - 例: v1.0.0, v1.0.1, v1.1.0, v2.0.0

2. **mainブランチではない**
   - `git checkout main` でmainブランチに切り替えてください

3. **未コミット変更がある**
   - `git add .` と `git commit -m "message"` でコミットしてください

4. **タグが既に存在**
   - 別のバージョン番号を使用するか、既存タグを削除してください

5. **GitHub Actionsが失敗**
   - ログを確認してエラー原因を特定してください
   - 必要に応じてコードを修正して再実行してください

6. **GitHub Actionsでリリース作成が失敗（403エラー）**
   - 権限不足でリリース作成に失敗する場合があります
   - 手動でリリースを作成してください：
   ```bash
   gh release create v0.1.2 --title "v0.1.2" --notes "リリースノート"
   ```
   - 注意: 手動リリース後もGitHub Actionsのビルドは継続されます

## ヘルパーコマンド

- `/update-homebrew` - Homebrew Formulaの更新
- `/test-homebrew` - Homebrewインストールのテスト
