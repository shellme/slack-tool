class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.2.1.tar.gz"
  sha256 "f99e13d3867453dd27a85a27f7f9c463263eb3bdb1361a504c2635e7f4e0e755"
  license "MIT"

  depends_on "go" => :build

  def install
    # バージョン情報を取得
    version_tag = "v#{version}"
    commit_hash = `git rev-parse --short HEAD`.strip
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
