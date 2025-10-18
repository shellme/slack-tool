#!/bin/bash

# SHA256計算ヘルパースクリプト
# 使用方法: ./calculate-sha256.sh [version]

set -e

# 色付き出力のための定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ関数
show_help() {
    echo "SHA256計算ヘルパースクリプト"
    echo ""
    echo "使用方法:"
    echo "  $0 [version]"
    echo ""
    echo "パラメータ:"
    echo "  version   対象のバージョン (例: v0.1.2)"
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
    
    # ファイルサイズを確認
    local file_size=$(stat -f%z "$temp_file" 2>/dev/null || stat -c%s "$temp_file" 2>/dev/null || echo "0")
    if [[ $file_size -eq 0 ]]; then
        log_error "ダウンロードしたファイルが空です"
        rm -f "$temp_file"
        exit 1
    fi
    
    log_info "ダウンロード完了: ${file_size} bytes"
    
    # SHA256を計算
    local sha256=$(shasum -a 256 "$temp_file" | cut -d' ' -f1)
    rm -f "$temp_file"
    
    log_success "SHA256計算完了: $sha256"
    echo "$sha256"
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
    
    log_info "SHA256計算開始: $version"
    echo ""
    
    # バージョン検証
    validate_version "$version"
    
    # SHA256計算
    local sha256=$(calculate_sha256 "$version")
    
    echo ""
    log_success "計算完了: $version -> $sha256"
    echo ""
    log_info "Formulaファイルの更新例:"
    echo "  url \"https://github.com/shellme/slack-tool/archive/$version.tar.gz\""
    echo "  sha256 \"$sha256\""
}

# スクリプト実行
main "$@"
