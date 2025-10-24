class Ko < Formula
  desc "Git Worktree tmux Automation - Manage git worktrees with tmux sessions"
  homepage "https://github.com/bshakr/ko"
  url "https://github.com/bshakr/ko.git",
      using:    :git,
      tag:      "v0.1.0",
      revision: "3bceb77e8c4c5a72d8821136fa0325a456078690"
  license "MIT"
  head "https://github.com/bshakr/ko.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the binary with version injection
    ldflags = "-s -w -X github.com/bshakr/ko/cmd.Version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "-o", bin/"ko"
  end

  test do
    # Test that the binary runs and shows version/help
    assert_match "Git Worktree tmux Automation", shell_output("#{bin}/ko --help")
  end
end
