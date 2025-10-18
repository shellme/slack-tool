class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.1.1.tar.gz"
  sha256 "833cf0dabff10dd0a054e69db412da75a6164ad523d12d447ab3a4526dd49b9a"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
