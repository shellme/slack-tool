# クイックリリースガイド

## 🚀 5ステップでリリース完了

```bash
# 1. 事前準備
git checkout main && git pull origin main

# 2. ローカル検証
/verify-build

# 3. リリース
/release v0.1.5 "リリースメッセージ"

# 4. Homebrew更新
/update-homebrew v0.1.5

# 5. テスト
/test-homebrew
```

## ⚡ 各コマンドの概要

| コマンド | 用途 | 時間 | 重要度 |
|---------|------|------|--------|
| `/verify-build` | ローカルビルド検証 | 1-2分 | ⭐⭐⭐⭐⭐ |
| `/release` | 新バージョンリリース | 3-5分 | ⭐⭐⭐⭐⭐ |
| `/update-homebrew` | Homebrew Formula更新 | 1-2分 | ⭐⭐⭐⭐⭐ |
| `/test-homebrew` | Homebrewインストールテスト | 2-3分 | ⭐⭐⭐⭐ |

## 🎯 成功の確認

```bash
# バージョン確認
slack-tool --version

# ヘルプ表示
slack-tool --help
```

## ⚠️ 注意事項

- **順序通りに実行**してください
- 前のコマンドが成功してから次を実行
- エラー時は各コマンドのドキュメントを確認

## 📚 詳細情報

- [完全なリリースワークフロー](./RELEASE_WORKFLOW.md)
- [各コマンドの詳細ドキュメント](../.cursor/commands/)
