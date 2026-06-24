# sample-service

Minimal Go HTTP service used to demonstrate the gitops-promoter promotion pipeline.

This repo contains **source code and CI only**. All delivery config — Argo CD Applications, gitops-promoter CRs, and the dry Helm chart — lives in [sample-service-config](https://github.com/platform-engineer-lab/sample-service-config).

## Endpoints

- `GET /` — returns `sample-service <version>`
- `GET /healthz` — returns `ok`

## Delivery pipeline

```mermaid
flowchart LR
    subgraph github["GitHub"]
        SRC["sample-service\nsource + Dockerfile + CI"]
        CFG["sample-service-config\ndry Helm chart"]
        ADD["platform-addons\nroles/&lt;role&gt;/"]
        APP["platform-apps\nregistry/*.yaml"]
    end

    GHCR[("GHCR\nghcr.io/…/sample-service:&lt;sha&gt;")]

    subgraph mgmt["management cluster"]
        AC(["Argo CD"])
        HY["Source\nHydrator"]
        GP["gitops-\npromoter"]
    end

    DEV(["dev spoke"])
    PROD(["prod spoke"])

    SRC -->|"CI: build + push :sha"| GHCR
    SRC -->|"CI: PR bump image.tag"| CFG
    ADD -->|"App-of-Apps"| AC
    APP -->|"cd-apps ApplicationSet"| AC
    CFG -->|"dry source HEAD"| HY
    HY -->|"push env/dev-next\nenv/prod-next"| CFG
    CFG -->|"env/dev · env/prod"| AC
    GP -->|"merge env/*-next → env/*"| CFG
    AC -->|"sync"| DEV
    AC -->|"sync"| PROD
    GHCR -.->|"pull"| DEV
    GHCR -.->|"pull"| PROD
    DEV -->|"argocd-health ✓\nunlocks prod"| GP
```

## How CI works

On every push to `main`:

1. **Test** — `go test ./...`
2. **Build + push** — multi-stage Docker build, image pushed to GHCR as
   `ghcr.io/platform-engineer-lab/sample-service:<short-sha>`
   (uses the built-in `GITHUB_TOKEN`; make the package public in GitHub settings
   so spoke clusters can pull without credentials)
3. **Bump tag** — opens a PR to `sample-service-config` updating `chart/values.yaml`
   `image.tag` to the new SHA. Once merged, the Source Hydrator re-renders the chart,
   gitops-promoter auto-merges the proposal to `env/dev`, and after dev reports healthy
   the prod promotion follows automatically.

## Required secrets

CI uses a GitHub App to open PRs into `sample-service-config`. Add two secrets in
this repo's Settings → Secrets → Actions:

| Secret | Value |
|---|---|
| `APP_ID` | GitHub App ID (e.g. `4117391`) |
| `APP_PRIVATE_KEY` | Contents of the downloaded `.pem` private key file |

The App must be installed on both `sample-service` and `sample-service-config` with
**Contents: read/write** and **Pull requests: read/write** permissions.

## Local dev

```bash
go test ./...
go run .
curl localhost:8080/healthz
```
