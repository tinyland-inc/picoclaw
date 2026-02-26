#!/bin/sh
# Substitute environment variables into PicoClaw config at startup.
# This runs before the PicoClaw gateway binary.

CONFIG="/home/picoclaw/.picoclaw/config.json"

if [ -f "$CONFIG" ]; then
  sed -i \
    -e "s|__ANTHROPIC_API_KEY__|${ANTHROPIC_API_KEY:-}|g" \
    -e "s|__ANTHROPIC_BASE_URL__|${ANTHROPIC_BASE_URL:-https://api.anthropic.com}|g" \
    "$CONFIG"
fi

exec picoclaw "$@"
