# Cursorカスタムコマンド

このディレクトリには、slack-toolのリリース作業を自動化するCursorカスタムコマンドが含まれています。

## 利用可能なコマンド

### `/release [version] [message]`
新バージョンのリリースを実行します。

**例:**
```
/release v0.1.2 "バグ修正とパフォーマンス改善"
/release v0.2.0
```

**機能:**
- バージョン番号の検証
- Gitタグの作成とプッシュ
- GitHub Actionsの実行監視
- リリース成功の確認

### `/update-homebrew [version]`
Homebrew Formulaを更新します。

**例:**
```
/update-homebrew v0.1.2
/update-homebrew  # 最新リリースを更新
```

**機能:**
- tarballのSHA256計算
- Formulaファイルの更新
- homebrew-slack-toolリポジトリへの反映

### `/test-homebrew`
Homebrewでのインストールをテストします。

**例:**
```
/test-homebrew
```

**機能:**
- 既存インストールのアンインストール
- 最新バージョンのインストール
- 動作テストの実行

## ファイル構成

```
.cursor/
├── commands/           # コマンドドキュメント
│   ├── release.md
│   ├── update-homebrew.md
│   └── test-homebrew.md
└── scripts/           # 実行スクリプト
    ├── release.sh
    ├── update-homebrew.sh
    ├── test-homebrew.sh
    └── helpers/
        ├── calculate-sha256.sh
        └── check-release.sh
```

## 使用方法

### AIアシスタント向け
コマンドを実行すると、対応するスクリプトが自動的に実行されます。

### 人間向け
直接スクリプトを実行することも可能です：

```bash
# リリース実行
./.cursor/scripts/release.sh v0.1.2 "バグ修正"

# Homebrew更新
./.cursor/scripts/update-homebrew.sh v0.1.2

# インストールテスト
./.cursor/scripts/test-homebrew.sh
```

## 前提条件

- GitHub CLI (gh) がインストールされている
- shellmeアカウントで認証されている
- Homebrewがインストールされている（テスト用）

## トラブルシューティング

各コマンドのドキュメント（`commands/`ディレクトリ）に詳細なトラブルシューティング情報が記載されています。

## 開発・カスタマイズ

スクリプトは`scripts/`ディレクトリにあり、必要に応じてカスタマイズできます。変更後は実行権限を確認してください：

```bash
chmod +x .cursor/scripts/*.sh
chmod +x .cursor/scripts/helpers/*.sh
```
