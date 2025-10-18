#!/bin/bash

# リリース確認ヘルパースクリプト
# 使用方法: ./check-release.sh [version]

set -e

# 色付き出力のための定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ関数
show_help() {
    echo "リリース確認ヘルパースクリプト"
    echo ""
    echo "使用方法:"
    echo "  $0 [version]"
    echo ""
    echo "パラメータ:"
    echo "  version   確認対象のバージョン (例: v0.1.2)"
    echo ""
    echo "例:"
    echo "  $0 v0.1.2"
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
        exit 1
    fi
    
    if ! gh auth status &> /dev/null; then
        log_error "GitHub CLI で認証されていません"
        exit 1
    fi
}

# リリースの存在確認
check_release_exists() {
    local version=$1
    
    if gh release view "$version" --repo shellme/slack-tool &> /dev/null; then
        log_success "リリース $version が存在します"
        return 0
    else
        log_error "リリース $version が見つかりません"
        return 1
    fi
}

# GitHub Actionsの実行状態確認
check_github_actions() {
    local version=$1
    
    log_info "GitHub Actionsの実行状態を確認中..."
    
    local runs=$(gh run list --repo shellme/slack-tool --limit 10 --json status,conclusion,headBranch,createdAt)
    local target_run=$(echo "$runs" | jq -r ".[] | select(.headBranch == \"main\") | select(.createdAt | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime > (now - 3600)) | ." | head -n1)
    
    if [[ -z "$target_run" || "$target_run" == "null" ]]; then
        log_warning "最近のGitHub Actions実行が見つかりません"
        return 1
    fi
    
    local status=$(echo "$target_run" | jq -r '.status')
    local conclusion=$(echo "$target_run" | jq -r '.conclusion')
    
    if [[ $status == "completed" ]]; then
        if [[ $conclusion == "success" ]]; then
            log_success "GitHub Actions実行完了: 成功"
            return 0
        else
            log_error "GitHub Actions実行失敗: $conclusion"
            return 1
        fi
    elif [[ $status == "in_progress" || $status == "queued" ]]; then
        log_warning "GitHub Actions実行中: $status"
        return 2
    else
        log_warning "GitHub Actions状態: $status"
        return 1
    fi
}

# リリースの詳細表示
show_release_details() {
    local version=$1
    
    log_info "リリース $version の詳細:"
    
    local release_info=$(gh release view "$version" --repo shellme/slack-tool --json name,url,assets,createdAt)
    local name=$(echo "$release_info" | jq -r '.name')
    local url=$(echo "$release_info" | jq -r '.url')
    local created_at=$(echo "$release_info" | jq -r '.createdAt')
    local asset_count=$(echo "$release_info" | jq -r '.assets | length')
    
    echo "  名前: $name"
    echo "  URL: $url"
    echo "  作成日時: $created_at"
    echo "  アセット数: $asset_count"
    
    if [[ $asset_count -gt 0 ]]; then
        echo "  アセット一覧:"
        echo "$release_info" | jq -r '.assets[] | "    - " + .name + " (" + (.size | tostring) + " bytes)"'
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
    
    log_info "リリース確認開始: $version"
    echo ""
    
    # 事前チェック
    check_gh_cli
    
    # リリース存在確認
    if check_release_exists "$version"; then
        echo ""
        show_release_details "$version"
        echo ""
        
        # GitHub Actions確認
        local actions_status
        check_github_actions "$version"
        actions_status=$?
        
        echo ""
        if [[ $actions_status -eq 0 ]]; then
            log_success "リリース $version は正常に完了しています"
            exit 0
        elif [[ $actions_status -eq 2 ]]; then
            log_warning "GitHub Actionsが実行中です。しばらく待ってから再確認してください"
            exit 2
        else
            log_error "GitHub Actionsの実行に問題があります"
            exit 1
        fi
    else
        log_error "リリース $version の確認に失敗しました"
        exit 1
    fi
}

# スクリプト実行
main "$@"
