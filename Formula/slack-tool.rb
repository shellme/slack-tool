class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.1.4.tar.gz"
  sha256 "e60f14d394fc9e176176ce02a7c912e5f26c0187aff5de8e243bf439f1a5ed2d"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", "#{bin}/slack-tool", "./cmd/slack-tool"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
