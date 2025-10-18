# リリースワークフロー

## 概要

slack-toolのリリースからHomebrew配布までの完全なワークフローです。カスタムコマンドを使用して効率的にリリース作業を行います。

## 🚀 完全なリリースフロー

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   事前準備      │ -> │ ローカル検証    │ -> │   リリース      │
│ git checkout    │    │ /verify-build   │    │ /release        │
│ git pull        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       v
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Homebrewテスト  │ <- │ Homebrew更新    │ <- │ GitHub Actions  │
│ /test-homebrew  │    │ /update-homebrew│    │ 実行中...       │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 1. 事前準備

```bash
# 現在のブランチを確認
git branch

# mainブランチに切り替え（必要に応じて）
git checkout main

# 最新の変更をプル
git pull origin main
```

### 2. ローカルビルド検証

```bash
/verify-build
```

**目的**: リリース前にローカルでビルド・実行テストを行い、問題を早期発見

**実行内容**:
- Go環境の確認
- 依存関係の整理
- コードの構文チェック
- ビルドテスト
- 実行テスト（バージョン確認、ヘルプ表示）

**成功例**:
```
[SUCCESS] Go環境の確認完了: go1.21.0
[SUCCESS] 依存関係の整理完了
[SUCCESS] コードの構文チェック完了
[SUCCESS] ビルドテスト完了
[SUCCESS] 実行テスト完了
[SUCCESS] すべての検証が成功しました！
```

### 3. 新バージョンのリリース

```bash
/release v0.1.5 "新機能追加とバグ修正"
```

**目的**: 指定したバージョンでGitタグを作成し、GitHub Actionsを実行してリリースを自動化

**実行内容**:
- バージョン番号の検証
- 環境の確認（mainブランチ、未コミット変更なし）
- Gitタグの作成
- GitHubへのプッシュ
- GitHub Actionsの実行監視
- 必要に応じて手動リリース作成

**成功例**:
```
[SUCCESS] タグ v0.1.5 を作成しました
[SUCCESS] タグをリモートにプッシュしました
[SUCCESS] GitHub Actionsを開始しました
[SUCCESS] リリースが正常に作成されました
```

### 4. Homebrew Formulaの更新

```bash
/update-homebrew v0.1.5
```

**目的**: 指定したバージョンのtarballのSHA256を計算し、Formulaファイルを更新してhomebrew-slack-toolリポジトリに反映

**実行内容**:
- リリース確認
- SHA256計算
- Formulaファイル更新
- homebrew-slack-toolリポジトリに反映

**成功例**:
```
[SUCCESS] リリース v0.1.5 を確認しました
[SUCCESS] SHA256計算完了: a1b2c3d4e5f6...
[SUCCESS] Formulaファイルを更新しました
[SUCCESS] homebrew-slack-toolリポジトリを更新しました
```

### 5. Homebrewインストールテスト

```bash
/test-homebrew
```

**目的**: Homebrewでのslack-toolインストールをテストし、動作確認を行う

**実行内容**:
- 既存インストールのアンインストール
- tapの確認・追加
- キャッシュのクリア
- 新規インストール
- 動作テスト

**成功例**:
```
[SUCCESS] 既存のインストールをアンインストールしました
[SUCCESS] tapを追加しました
[SUCCESS] インストールが完了しました
[SUCCESS] バージョン確認: slack-tool version 1.0.0
[SUCCESS] ヘルプ表示: 正常
[SUCCESS] すべてのテストが成功しました！
```

## 📋 各コマンドの詳細

### `/verify-build`
- **用途**: リリース前のローカル検証
- **実行時間**: 1-2分
- **重要度**: ⭐⭐⭐⭐⭐（必須）

### `/release`
- **用途**: 新バージョンのリリース
- **実行時間**: 3-5分（GitHub Actions含む）
- **重要度**: ⭐⭐⭐⭐⭐（必須）

### `/update-homebrew`
- **用途**: Homebrew Formulaの更新
- **実行時間**: 1-2分
- **重要度**: ⭐⭐⭐⭐⭐（必須）

### `/test-homebrew`
- **用途**: Homebrewインストールのテスト
- **実行時間**: 2-3分
- **重要度**: ⭐⭐⭐⭐（推奨）

## 🔄 完全なリリース例

```bash
# 1. 事前準備
git checkout main
git pull origin main

# 2. ローカルビルド検証
/verify-build

# 3. 新バージョンのリリース
/release v0.1.5 "新機能追加とバグ修正"

# 4. Homebrew Formulaの更新
/update-homebrew v0.1.5

# 5. Homebrewインストールテスト
/test-homebrew
```

## ⚠️ 注意事項

### 実行順序
1. **必ず順序通りに実行**してください
2. 前のコマンドが成功してから次のコマンドを実行してください
3. エラーが発生した場合は、エラーを解決してから再実行してください

### 前提条件
- GitHub CLI (gh) がインストールされている
- shellmeアカウントで認証されている
- リモートリポジトリが設定されている
- プッシュ権限がある
- Homebrewがインストールされている

### トラブルシューティング
- 各コマンドの詳細なトラブルシューティング情報は、各コマンドのドキュメントを参照してください
- 問題が解決しない場合は、手動でコマンドを実行してください

## 🎯 成功の確認

すべてのコマンドが成功した場合、以下が確認できます：

1. **GitHubリリース**: https://github.com/shellme/slack-tool/releases
2. **Homebrewインストール**: `brew install shellme/slack-tool/slack-tool`
3. **動作確認**: `slack-tool --version`

## 📚 関連ドキュメント

- [リリースコマンド](./release.md)
- [Homebrew更新コマンド](./update-homebrew.md)
- [Homebrewテストコマンド](./test-homebrew.md)
- [ビルド検証コマンド](./verify-build.md)

## 🚀 次のステップ

リリースが完了したら：

1. **ユーザーへの通知**: 新バージョンのリリースをユーザーに通知
2. **ドキュメント更新**: 必要に応じてREADMEやドキュメントを更新
3. **次の開発**: 次のバージョンの開発を開始

---

**注意**: このワークフローは、slack-toolのリリースプロセスに特化しています。他のプロジェクトでは適宜調整してください。
