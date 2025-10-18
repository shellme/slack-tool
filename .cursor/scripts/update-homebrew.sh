#!/bin/bash

# Homebrew Formula更新スクリプト
# 使用方法: ./update-homebrew.sh [version]

set -e

# 色付き出力のための定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ関数
show_help() {
    echo "Homebrew Formula更新スクリプト"
    echo ""
    echo "使用方法:"
    echo "  $0 [version]"
    echo ""
    echo "パラメータ:"
    echo "  version   更新対象のバージョン (省略時は最新リリース)"
    echo ""
    echo "例:"
    echo "  $0 v0.1.2"
    echo "  $0"
    echo ""
    echo "オプション:"
    echo "  --help    このヘルプを表示"
}

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# GitHub CLIの確認
check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) がインストールされていません"
        log_info "以下のコマンドでインストールしてください:"
        echo "  brew install gh"
        exit 1
    fi
    
    if ! gh auth status &> /dev/null; then
        log_error "GitHub CLI で認証されていません"
        log_info "以下のコマンドで認証してください:"
        echo "  gh auth login"
        exit 1
    fi
    
    log_success "GitHub CLI確認完了"
}

# 最新リリースの取得
get_latest_release() {
    local latest_release=$(gh release list --repo shellme/slack-tool --limit 1 --json tagName -q '.[0].tagName')
    if [[ -z "$latest_release" || "$latest_release" == "null" ]]; then
        log_error "リリースが見つかりません"
        log_info "先に /release コマンドでリリースを作成してください"
        exit 1
    fi
    echo "$latest_release"
}

# リリースの存在確認
check_release_exists() {
    local version=$1
    if ! gh release view "$version" --repo shellme/slack-tool &> /dev/null; then
        log_error "リリース $version が見つかりません"
        log_info "利用可能なリリース:"
        gh release list --repo shellme/slack-tool --limit 10
        exit 1
    fi
    log_success "リリース確認完了: $version"
}

# SHA256の計算
calculate_sha256() {
    local version=$1
    local tarball_url="https://github.com/shellme/slack-tool/archive/$version.tar.gz"
    
    log_info "tarballをダウンロード中: $tarball_url"
    
    # 一時ファイルにダウンロード
    local temp_file=$(mktemp)
    if ! curl -L -s "$tarball_url" -o "$temp_file"; then
        log_error "tarballのダウンロードに失敗しました"
        rm -f "$temp_file"
        exit 1
    fi
    
    # SHA256を計算
    local sha256=$(shasum -a 256 "$temp_file" | cut -d' ' -f1)
    rm -f "$temp_file"
    
    log_success "SHA256計算完了: $sha256"
    echo "$sha256"
}

# Formulaファイルの更新
update_formula() {
    local version=$1
    local sha256=$2
    local formula_file="Formula/slack-tool.rb"
    
    if [[ ! -f "$formula_file" ]]; then
        log_error "Formulaファイルが見つかりません: $formula_file"
        exit 1
    fi
    
    log_info "Formulaファイルを更新中: $formula_file"
    
    # バックアップを作成
    cp "$formula_file" "$formula_file.backup"
    
    # URLとSHA256を更新
    local tarball_url="https://github.com/shellme/slack-tool/archive/$version.tar.gz"
    
    # sedを使用してURLとSHA256を更新
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS用
        sed -i '' "s|url \".*\"|url \"$tarball_url\"|g" "$formula_file"
        sed -i '' "s|sha256 \".*\"|sha256 \"$sha256\"|g" "$formula_file"
    else
        # Linux用
        sed -i "s|url \".*\"|url \"$tarball_url\"|g" "$formula_file"
        sed -i "s|sha256 \".*\"|sha256 \"$sha256\"|g" "$formula_file"
    fi
    
    log_success "Formulaファイル更新完了"
    
    # 変更内容を表示
    log_info "更新内容:"
    echo "  URL: $tarball_url"
    echo "  SHA256: $sha256"
}

# homebrew-slack-toolリポジトリの更新
update_homebrew_repo() {
    local version=$1
    local temp_dir="/tmp/homebrew-slack-tool-$$"
    
    log_info "homebrew-slack-toolリポジトリを更新中..."
    
    # 既存のディレクトリを削除
    rm -rf "$temp_dir"
    
    # リポジトリをクローン
    if ! git clone https://github.com/shellme/homebrew-slack-tool.git "$temp_dir"; then
        log_error "homebrew-slack-toolリポジトリのクローンに失敗しました"
        log_info "リポジトリが存在するか確認してください:"
        echo "  https://github.com/shellme/homebrew-slack-tool"
        exit 1
    fi
    
    cd "$temp_dir"
    
    # Formulaファイルをコピー
    if ! cp "$OLDPWD/Formula/slack-tool.rb" "./Formula/slack-tool.rb"; then
        log_error "Formulaファイルのコピーに失敗しました"
        cd "$OLDPWD"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # 変更をコミット
    git add Formula/slack-tool.rb
    if git diff --staged --quiet; then
        log_warning "変更がありません（既に最新の状態です）"
        cd "$OLDPWD"
        rm -rf "$temp_dir"
        return 0
    fi
    
    git commit -m "Update slack-tool to $version"
    
    # プッシュ
    if ! git push origin main; then
        log_error "プッシュに失敗しました"
        log_info "権限を確認してください"
        cd "$OLDPWD"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    log_success "homebrew-slack-toolリポジトリ更新完了"
    
    # クリーンアップ
    cd "$OLDPWD"
    rm -rf "$temp_dir"
}

# Formulaファイルの検証
validate_formula() {
    local formula_file="Formula/slack-tool.rb"
    
    log_info "Formulaファイルの検証中..."
    
    # 構文チェック
    if ! ruby -c "$formula_file" &> /dev/null; then
        log_error "Formulaファイルの構文エラーがあります"
        log_info "以下のコマンドで詳細を確認してください:"
        echo "  ruby -c $formula_file"
        exit 1
    fi
    
    # Homebrewのauditチェック（可能な場合）
    if command -v brew &> /dev/null; then
        log_info "Homebrew auditを実行中..."
        if brew audit --strict "$formula_file" &> /dev/null; then
            log_success "Homebrew audit完了: 問題なし"
        else
            log_warning "Homebrew auditで警告があります"
            log_info "詳細を確認してください:"
            echo "  brew audit --strict $formula_file"
        fi
    else
        log_warning "Homebrewがインストールされていないため、auditをスキップします"
    fi
}

# メイン処理
main() {
    # 引数チェック
    if [[ $# -gt 0 && ($1 == "--help" || $1 == "-h") ]]; then
        show_help
        exit 0
    fi
    
    local version=${1:-""}
    
    log_info "Homebrew Formula更新開始"
    echo ""
    
    # 事前チェック
    check_gh_cli
    
    # バージョン決定
    if [[ -z "$version" ]]; then
        log_info "最新リリースを取得中..."
        version=$(get_latest_release)
        log_info "対象バージョン: $version"
    else
        check_release_exists "$version"
    fi
    
    echo ""
    log_info "更新対象: slack-tool $version"
    echo ""
    
    # SHA256計算
    local sha256=$(calculate_sha256 "$version")
    echo ""
    
    # Formulaファイル更新
    update_formula "$version" "$sha256"
    echo ""
    
    # Formulaファイル検証
    validate_formula
    echo ""
    
    # homebrew-slack-toolリポジトリ更新
    update_homebrew_repo "$version"
    echo ""
    
    log_success "Homebrew Formula更新完了: $version"
    log_info "次のステップ:"
    echo "  /test-homebrew  # インストールテストを実行"
    echo ""
    log_info "インストール方法:"
    echo "  brew tap shellme/slack-tool"
    echo "  brew install slack-tool"
    echo "  slack-tool --version"
}

# スクリプト実行
main "$@"
