#!/bin/sh
# Substitute environment variables into TinyClaw config at startup.
# This runs before the TinyClaw gateway binary.

CONFIG="/home/tinyclaw/.tinyclaw/config.json"

if [ -f "$CONFIG" ]; then
  sed -i \
    -e "s|__ANTHROPIC_API_KEY__|${ANTHROPIC_API_KEY:-}|g" \
    -e "s|__ANTHROPIC_BASE_URL__|${ANTHROPIC_BASE_URL:-https://api.anthropic.com}|g" \
    "$CONFIG"
fi

exec tinyclaw "$@"
