#!/bin/bash

# slack-tool リリーススクリプト
# 使用方法: ./release.sh [version] [message]

set -e

# 色付き出力のための定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ関数
show_help() {
    echo "slack-tool リリーススクリプト"
    echo ""
    echo "使用方法:"
    echo "  $0 [version] [message]"
    echo ""
    echo "パラメータ:"
    echo "  version   リリースバージョン (例: v0.1.2, v0.2.0)"
    echo "  message   リリースメッセージ (省略可)"
    echo ""
    echo "例:"
    echo "  $0 v0.1.2 \"バグ修正とパフォーマンス改善\""
    echo "  $0 v0.2.0"
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

# バージョン番号の検証
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "バージョン番号が不正です: $version"
        log_error "セマンティックバージョニング形式を使用してください (例: v1.0.0)"
        exit 1
    fi
}

# 現在のブランチ確認
check_branch() {
    local current_branch=$(git branch --show-current)
    if [[ $current_branch != "main" ]]; then
        log_error "mainブランチで実行してください。現在のブランチ: $current_branch"
        log_info "以下のコマンドでmainブランチに切り替えてください:"
        echo "  git checkout main"
        exit 1
    fi
    log_success "ブランチ確認完了: $current_branch"
}

# 未コミット変更の確認
check_uncommitted_changes() {
    if ! git diff-index --quiet HEAD --; then
        log_error "未コミットの変更があります"
        log_info "以下のコマンドで変更をコミットしてください:"
        echo "  git add ."
        echo "  git commit -m \"your message\""
        exit 1
    fi
    log_success "未コミット変更の確認完了"
}

# タグの重複確認
check_tag_exists() {
    local version=$1
    if git tag -l | grep -q "^$version$"; then
        log_error "タグ $version は既に存在します"
        log_info "別のバージョン番号を使用するか、既存タグを削除してください:"
        echo "  git tag -d $version"
        echo "  git push origin :refs/tags/$version"
        exit 1
    fi
    log_success "タグ重複確認完了: $version"
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

# リモートリポジトリの確認
check_remote() {
    if ! git remote get-url origin &> /dev/null; then
        log_error "リモートリポジトリが設定されていません"
        log_info "以下のコマンドでリモートリポジトリを設定してください:"
        echo "  git remote add origin https://github.com/shellme/slack-tool.git"
        exit 1
    fi
    
    local remote_url=$(git remote get-url origin)
    log_success "リモートリポジトリ確認完了: $remote_url"
}

# タグの作成とプッシュ
create_and_push_tag() {
    local version=$1
    local message=$2
    
    log_info "タグ $version を作成中..."
    git tag -a "$version" -m "$message"
    log_success "タグ作成完了: $version"
    
    log_info "タグをプッシュ中..."
    git push origin "$version"
    log_success "タグプッシュ完了: $version"
}

# GitHub Actionsの実行監視
monitor_github_actions() {
    local version=$1
    
    log_info "GitHub Actionsの実行を監視中..."
    
    # 最大10分間待機
    local max_wait=600
    local wait_time=0
    local check_interval=30
    
    while [ $wait_time -lt $max_wait ]; do
        local runs=$(gh run list --repo shellme/slack-tool --limit 5 --json status,conclusion,headBranch)
        local latest_run=$(echo "$runs" | jq -r '.[0]')
        local status=$(echo "$latest_run" | jq -r '.status')
        local conclusion=$(echo "$latest_run" | jq -r '.conclusion')
        local branch=$(echo "$latest_run" | jq -r '.headBranch')
        
        if [[ $branch == "main" ]]; then
            if [[ $status == "completed" ]]; then
                if [[ $conclusion == "success" ]]; then
                    log_success "GitHub Actions実行完了: 成功"
                    return 0
                else
                    log_error "GitHub Actions実行失敗: $conclusion"
                    log_info "詳細なログを確認してください:"
                    echo "  gh run list --repo shellme/slack-tool"
                    echo "  gh run view <run-id> --log"
                    return 1
                fi
            elif [[ $status == "in_progress" || $status == "queued" ]]; then
                log_info "GitHub Actions実行中... ($status) - ${wait_time}秒経過"
            else
                log_warning "GitHub Actions状態: $status"
            fi
        fi
        
        sleep $check_interval
        wait_time=$((wait_time + check_interval))
    done
    
    log_warning "GitHub Actionsの監視がタイムアウトしました"
    log_info "手動で確認してください:"
    echo "  gh run list --repo shellme/slack-tool"
    return 1
}

# リリース確認
verify_release() {
    local version=$1
    
    log_info "リリース $version の確認中..."
    
    # リリースが作成されているか確認
    if gh release view "$version" --repo shellme/slack-tool &> /dev/null; then
        log_success "リリース $version が正常に作成されました"
        
        # リリースの詳細を表示
        log_info "リリース詳細:"
        gh release view "$version" --repo shellme/slack-tool --json name,url,assets | jq -r '
            "  名前: " + .name + "\n" +
            "  URL: " + .url + "\n" +
            "  アセット数: " + (.assets | length | tostring)
        '
        
        return 0
    else
        log_error "リリース $version が見つかりません"
        return 1
    fi
}

# メイン処理
main() {
    # 引数チェック
    if [[ $# -eq 0 || $1 == "--help" || $1 == "-h" ]]; then
        show_help
        exit 0
    fi
    
    if [[ $# -lt 1 ]]; then
        log_error "バージョン番号を指定してください"
        show_help
        exit 1
    fi
    
    local version=$1
    local message=${2:-"Release $version"}
    
    log_info "slack-tool リリース開始: $version"
    log_info "リリースメッセージ: $message"
    echo ""
    
    # 事前チェック
    validate_version "$version"
    check_gh_cli
    check_remote
    check_branch
    check_uncommitted_changes
    check_tag_exists "$version"
    
    echo ""
    log_info "すべての事前チェックが完了しました"
    echo ""
    
    # リリース実行
    create_and_push_tag "$version" "$message"
    echo ""
    
    # GitHub Actions監視
    if monitor_github_actions "$version"; then
        echo ""
        # リリース確認
        if verify_release "$version"; then
            echo ""
            log_success "リリース $version が正常に完了しました！"
            log_info "次のステップ:"
            echo "  /update-homebrew $version  # Homebrew Formulaを更新"
            echo "  /test-homebrew            # インストールテスト"
        else
            log_error "リリース確認に失敗しました"
            exit 1
        fi
    else
        log_error "GitHub Actionsの実行に失敗しました"
        exit 1
    fi
}

# スクリプト実行
main "$@"
