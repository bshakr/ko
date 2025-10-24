class Ko < Formula
  desc "Git Worktree tmux Automation - Manage git worktrees with tmux sessions"
  homepage "https://github.com/bshakr/ko"
  url "https://github.com/bshakr/ko.git",
      using:    :git,
      tag:      "v0.1.1",
      revision: "21fbb1a30e0b775a02b75b2c1e73fea1c10bf3d1"
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
