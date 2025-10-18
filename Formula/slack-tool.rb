class SlackTool < Formula
  desc "Slackの様々な操作を行うCLIツール"
  homepage "https://github.com/shellme/slack-tool"
  url "https://github.com/shellme/slack-tool/archive/v0.1.5.tar.gz"
  sha256 "b97d69b7129e9324586df8ff1c0188f7fde68e9e93c5e72203c8020a30a9066e"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", "#{bin}/slack-tool", "./cmd/slack-tool"
  end

  test do
    assert_match "slack-tool version", shell_output("#{bin}/slack-tool --version")
  end
end
