# typed: false
# frozen_string_literal: true

class Tinyclaw < Formula
  desc "Ultra-lightweight personal AI agent"
  homepage "https://github.com/tinyland-inc/tinyclaw"
  license "MIT"
  # version/url/sha256 filled by goreleaser for bottled releases
  head "https://github.com/tinyland-inc/tinyclaw.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal.version=#{version}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags), "-tags", "stdjson", "./cmd/tinyclaw"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/tinyclaw version")
  end
end
