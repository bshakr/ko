class Ko < Formula
  desc "Git worktree + tmux automation for isolated development environments"
  homepage "https://github.com/bshakr/ko"
  url "https://github.com/bshakr/ko/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # Will be filled after creating the release
  license "MIT"

  def install
    bin.install "ko.sh" => "ko"
  end

  test do
    assert_match "Usage: ko", shell_output("#{bin}/ko --help", 1)
  end
end
