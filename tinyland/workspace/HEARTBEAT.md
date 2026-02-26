# PicoClaw Heartbeat Checklist

Runs every 120 minutes during active hours (06:00-22:00).

## Quick Health Check

1. **Workspace integrity**: Verify AGENT.md, IDENTITY.md, SOUL.md exist
2. **Memory health**: Check MEMORY.md is accessible and not empty
3. **Tool availability**: Verify adapter tool proxy is responsive
4. **Upstream status**: Note if upstream check is overdue

## Periodic Tasks

### Weekly (Monday): Upstream Sync Check
- Fetch sipeed/picoclaw main HEAD via campaign
- Compare with tinyland-inc/picoclaw tinyland HEAD
- Log delta to MEMORY.md

### Daily: Memory Compaction
- Remove duplicate observations
- Consolidate recurring patterns
- Prune entries older than 14 days

## Escalation

- If tools unreachable for 2+ consecutive heartbeats: log and wait
- If workspace files missing: restore from /workspace-defaults/
- If memory corrupted: reset from seed and log event
