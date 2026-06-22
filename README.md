# sample-service

Minimal Go HTTP service used to demonstrate the gitops-promoter promotion pipeline.

## Endpoints

- `GET /` — returns `sample-service <version>`
- `GET /healthz` — returns `ok`

## How CI works

On every push to `main`:

1. **Test** — `go test ./...`
2. **Build + push** — multi-stage Docker build, image pushed to GHCR as
   `ghcr.io/platform-engineer-lab/sample-service:<short-sha>`
   (uses the built-in `GITHUB_TOKEN`; make the package public in GitHub settings
   so spoke clusters can pull without credentials)
3. **Bump tag** — opens a PR to `sample-service-config` updating `chart/values.yaml`
   `image.tag` to the new SHA.

## Required secret

Add a fine-grained PAT as `CONFIG_REPO_TOKEN` in this repo's Settings → Secrets:
- **Repository access:** `platform-engineer-lab/sample-service-config`
- **Permissions:** Contents (read/write), Pull requests (read/write)

## Local dev

```bash
go test ./...
go run .
curl localhost:8080/healthz
```
