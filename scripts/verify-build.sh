#!/bin/bash

# ビルド検証スクリプト
# slack-toolのローカルビルドを検証します

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# ヘルプ表示
show_help() {
    cat << EOF
ビルド検証スクリプト

使用方法:
    $0 [オプション]

オプション:
    -h, --help     このヘルプを表示
    -v, --verbose  詳細な出力を表示

説明:
    slack-toolのローカルビルドを検証します。
    リリース前に必ず実行してください。

EOF
}

# 引数解析
VERBOSE=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        *)
            log_error "不明なオプション: $1"
            show_help
            exit 1
            ;;
    esac
done

# メイン処理
main() {
    log_info "slack-toolのビルド検証を開始します"
    
    # 1. 環境確認
    log_info "Go環境の確認中..."
    if ! command -v go &> /dev/null; then
        log_error "Goがインストールされていません"
        exit 1
    fi
    
    go_version=$(go version | cut -d' ' -f3)
    log_success "Go環境の確認完了: $go_version"
    
    # 2. 依存関係の整理
    log_info "依存関係の整理中..."
    if ! go mod tidy; then
        log_error "依存関係の整理に失敗しました"
        exit 1
    fi
    log_success "依存関係の整理完了"
    
    # 3. コードの構文チェック
    log_info "コードの構文チェック中..."
    if ! go vet ./...; then
        log_error "コードの構文チェックに失敗しました"
        exit 1
    fi
    log_success "コードの構文チェック完了"
    
    # 4. コードフォーマットの確認
    log_info "コードフォーマットの確認中..."
    if ! go fmt ./...; then
        log_warning "コードフォーマットに問題があります"
    else
        log_success "コードフォーマットの確認完了"
    fi
    
    # 5. ビルドテスト
    log_info "ビルドテスト中..."
    if ! go build -o slack-tool ./cmd/slack-tool; then
        log_error "ビルドに失敗しました"
        exit 1
    fi
    log_success "ビルドテスト完了"
    
    # 6. 実行権限の設定
    log_info "実行権限を設定中..."
    if ! chmod +x slack-tool; then
        log_error "実行権限の設定に失敗しました"
        exit 1
    fi
    log_success "実行権限の設定完了"
    
    # 7. 実行テスト
    log_info "実行テスト中..."
    
    # バージョン確認
    if ! ./slack-tool --version > /dev/null 2>&1; then
        log_error "バージョン確認に失敗しました"
        exit 1
    fi
    log_success "バージョン確認完了"
    
    # ヘルプ表示確認
    if ! ./slack-tool --help > /dev/null 2>&1; then
        log_error "ヘルプ表示に失敗しました"
        exit 1
    fi
    log_success "ヘルプ表示確認完了"
    
    # 8. クリーンアップ
    log_info "クリーンアップ中..."
    if ! rm -f slack-tool; then
        log_warning "テスト用ファイルの削除に失敗しました"
    else
        log_success "クリーンアップ完了"
    fi
    
    # 完了
    log_success "すべての検証が成功しました！"
    log_info "リリースの準備が整いました"
}

# スクリプト実行
main "$@"
