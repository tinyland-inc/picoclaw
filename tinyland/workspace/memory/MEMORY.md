# PicoClaw Long-Term Memory

## Platform Architecture

- **Cluster**: Civo Kubernetes, namespace `fuzzy-dev`
- **Gateway**: `http://rj-gateway.fuzzy-dev.svc.cluster.local:8080`
- **Aperture**: LLM proxy at `http://aperture.fuzzy-dev.svc.cluster.local`
- **Bot**: `rj-agent-bot[bot]` (GitHub App ID 2945224)

## Fork Status

- **Upstream**: sipeed/picoclaw (main branch)
- **Our fork**: tinyland-inc/picoclaw (tinyland branch)
- **Customizations**: Dockerfile, config.json, entrypoint.sh, workspace files
- **Last upstream check**: (not yet performed)
- **Commits behind upstream**: (unknown)

## Known False Positives

(Track patterns that consistently produce false positives to skip in future scans)

## Observations

(Populated by campaign results and heartbeat observations)

## Known Issues

(Recurring problems and their workarounds)
