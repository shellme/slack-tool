class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.2.0.tar.gz"
  sha256 "b5c2f065800fdf4b44a3c2895802d5de1b8eac9238cc322323b724f5e9748943"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", "#{bin}/slack-tool", "./cmd/slack-tool"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
