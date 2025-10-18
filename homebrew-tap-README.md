# Homebrew Tap for slack-tool

このリポジトリは [slack-tool](https://github.com/shellme/slack-tool) のHomebrew Formulaを提供します。

## インストール

```bash
# tapを追加
brew tap shellme/slack-tool

# slack-toolをインストール
brew install slack-tool
```

## アンインストール

```bash
# slack-toolをアンインストール
brew uninstall slack-tool

# tapを削除（オプション）
brew untap shellme/slack-tool
```

## 更新

```bash
# slack-toolを更新
brew upgrade slack-tool
```

## トラブルシューティング

### インストールに失敗する場合

```bash
# tapを更新してから再試行
brew update
brew tap shellme/slack-tool
brew install slack-tool
```

### 古いバージョンがインストールされる場合

```bash
# キャッシュをクリアしてから再インストール
brew uninstall slack-tool
brew untap shellme/slack-tool
brew tap shellme/slack-tool
brew install slack-tool
```
