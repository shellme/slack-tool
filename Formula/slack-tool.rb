class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.4.0.tar.gz"
  sha256 "d106b8ee9901df59386772900582ae6996c42ca3b4e3c08ceb1f11089584c81e"
  license "MIT"

  depends_on "go" => :build

  def install
    # バージョン情報を取得
    version_tag = "v#{version}"
    # tarballにgit情報がないためタグのコミットを指定
    commit_hash = "ab0aa77"
    build_date = `date -u '+%Y-%m-%d_%H:%M:%S'`.strip
    
    # ビルドフラグを設定
    ldflags = [
      "-s -w",
      "-X github.com/shellme/slack-tool/cmd/slack-tool/cmd.version=#{version_tag}",
      "-X github.com/shellme/slack-tool/cmd/slack-tool/cmd.commit=#{commit_hash}",
      "-X github.com/shellme/slack-tool/cmd/slack-tool/cmd.date=#{build_date}"
    ].join(" ")
    
    system "go", "build", "-ldflags", ldflags, "-o", "#{bin}/slack-tool", "./cmd/slack-tool"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
