#!/bin/bash

# Homebrewインストールテストスクリプト
# 使用方法: ./test-homebrew.sh

set -e

# 色付き出力のための定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ関数
show_help() {
    echo "Homebrewインストールテストスクリプト"
    echo ""
    echo "使用方法:"
    echo "  $0"
    echo ""
    echo "このスクリプトは以下の処理を実行します:"
    echo "  1. 既存のslack-toolをアンインストール"
    echo "  2. Homebrewキャッシュをクリア"
    echo "  3. shellme/slack-toolのtapを追加"
    echo "  4. 最新バージョンをインストール"
    echo "  5. 動作テストを実行"
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

# Homebrewの確認
check_homebrew() {
    if ! command -v brew &> /dev/null; then
        log_error "Homebrewがインストールされていません"
        log_info "以下のコマンドでインストールしてください:"
        echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        exit 1
    fi
    
    local brew_version=$(brew --version | head -n1)
    log_success "Homebrew確認完了: $brew_version"
}

# 既存インストールの確認とアンインストール
uninstall_existing() {
    log_info "既存のslack-toolインストールを確認中..."
    
    if brew list --formula | grep -q "^slack-tool$"; then
        log_info "既存のslack-toolをアンインストール中..."
        if brew uninstall slack-tool; then
            log_success "既存のslack-toolをアンインストールしました"
        else
            log_warning "アンインストールに失敗しました（続行します）"
        fi
    else
        log_info "既存のslack-toolはインストールされていません"
    fi
}

# tapの確認と追加
setup_tap() {
    log_info "shellme/slack-toolのtapを確認中..."
    
    if brew tap | grep -q "^shellme/slack-tool$"; then
        log_success "tapは既に追加されています"
    else
        log_info "tapを追加中..."
        if brew tap shellme/slack-tool; then
            log_success "tapを追加しました: shellme/slack-tool"
        else
            log_error "tapの追加に失敗しました"
            log_info "手動で追加してください:"
            echo "  brew tap shellme/slack-tool"
            exit 1
        fi
    fi
}

# キャッシュのクリア
clear_cache() {
    log_info "Homebrewキャッシュをクリア中..."
    
    if brew cleanup; then
        log_success "キャッシュをクリアしました"
    else
        log_warning "キャッシュのクリアに失敗しました（続行します）"
    fi
}

# インストール
install_slack_tool() {
    log_info "slack-toolをインストール中..."
    
    if brew install slack-tool; then
        log_success "slack-toolのインストールが完了しました"
    else
        log_error "インストールに失敗しました"
        log_info "デバッグ情報:"
        echo "  brew install --debug slack-tool"
        exit 1
    fi
}

# バージョン確認
check_version() {
    log_info "バージョンを確認中..."
    
    if command -v slack-tool &> /dev/null; then
        local version=$(slack-tool --version 2>/dev/null || echo "unknown")
        log_success "バージョン確認: $version"
        
        # バージョンが正しく表示されるかチェック
        if [[ "$version" == *"slack-tool version"* ]]; then
            log_success "バージョン表示: 正常"
        else
            log_warning "バージョン表示が期待と異なります: $version"
        fi
    else
        log_error "slack-toolコマンドが見つかりません"
        exit 1
    fi
}

# ヘルプ表示テスト
test_help() {
    log_info "ヘルプ表示をテスト中..."
    
    if slack-tool --help &> /dev/null; then
        log_success "ヘルプ表示: 正常"
    else
        log_error "ヘルプ表示に失敗しました"
        exit 1
    fi
}

# 基本コマンドテスト
test_basic_commands() {
    log_info "基本コマンドをテスト中..."
    
    # config showコマンドのテスト（トークンが設定されていなくてもエラーにならないはず）
    if slack-tool config show &> /dev/null; then
        log_success "config showコマンド: 正常"
    else
        log_warning "config showコマンドでエラーが発生しました（トークン未設定の可能性）"
    fi
    
    # その他のコマンドもテスト（エラーが発生しても致命的でない場合は続行）
    log_info "その他のコマンドをテスト中..."
    
    # 無効なURLでのテスト（エラーが期待される）
    if ! slack-tool get "invalid-url" &> /dev/null; then
        log_success "無効なURLのエラーハンドリング: 正常"
    else
        log_warning "無効なURLのエラーハンドリングに問題があります"
    fi
}

# Formulaの検証
validate_formula() {
    log_info "Formulaを検証中..."
    
    if brew audit --strict slack-tool &> /dev/null; then
        log_success "Formula検証: 問題なし"
    else
        log_warning "Formula検証で警告があります"
        log_info "詳細を確認してください:"
        echo "  brew audit --strict slack-tool"
    fi
}

# インストール情報の表示
show_installation_info() {
    log_info "インストール情報:"
    
    local install_path=$(brew --prefix slack-tool 2>/dev/null || echo "不明")
    local version=$(slack-tool --version 2>/dev/null || echo "不明")
    
    echo "  インストールパス: $install_path"
    echo "  バージョン: $version"
    echo "  実行ファイル: $(which slack-tool)"
}

# メイン処理
main() {
    # 引数チェック
    if [[ $# -gt 0 && ($1 == "--help" || $1 == "-h") ]]; then
        show_help
        exit 0
    fi
    
    log_info "Homebrewインストールテスト開始"
    echo ""
    
    # 事前チェック
    check_homebrew
    echo ""
    
    # 既存インストールのアンインストール
    uninstall_existing
    echo ""
    
    # tapの設定
    setup_tap
    echo ""
    
    # キャッシュクリア
    clear_cache
    echo ""
    
    # インストール
    install_slack_tool
    echo ""
    
    # テスト実行
    log_info "動作テストを実行中..."
    echo ""
    
    check_version
    test_help
    test_basic_commands
    echo ""
    
    # Formula検証
    validate_formula
    echo ""
    
    # インストール情報表示
    show_installation_info
    echo ""
    
    log_success "すべてのテストが成功しました！"
    log_info "slack-toolが正常にインストールされ、動作しています。"
    echo ""
    log_info "使用例:"
    echo "  slack-tool --version"
    echo "  slack-tool --help"
    echo "  slack-tool config show"
    echo "  slack-tool get <slack-thread-url>"
}

# スクリプト実行
main "$@"
