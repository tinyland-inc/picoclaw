//go:build nochinese

package channels

import (
	"github.com/tinyland-inc/tinyclaw/pkg/bus"
	"github.com/tinyland-inc/tinyclaw/pkg/config"
)

// initChineseChannels is a no-op stub when building with -tags nochinese.
// Chinese-only channels (Feishu, QQ, DingTalk, OneBot, WeCom, WeComApp)
// are excluded from the build.
func initChineseChannels(_ *config.Config, _ *bus.MessageBus, _ map[string]Channel) {
	// No Chinese channels available in this build.
}
