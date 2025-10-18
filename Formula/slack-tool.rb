class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/slack-tool"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
