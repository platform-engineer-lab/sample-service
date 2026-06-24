# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Delivery pipeline

```mermaid
flowchart LR
    subgraph github["GitHub"]
        SRC["sample-service\nsource + Dockerfile + CI"]
        CFG["sample-service-config\ndry Helm chart"]
        ADD["platform-addons\nroles/<role>/"]
        APP["platform-apps\nregistry/*.yaml"]
    end

    GHCR[("GHCR\nghcr.io/…/sample-service:<sha>")]

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

## Architecture

This repo is the source side of the pipeline. On every push to `main`, CI:

1. Runs `go test ./...`
2. Builds a multi-stage Docker image and pushes `ghcr.io/platform-engineer-lab/sample-service:<short-sha>` to GHCR using the built-in `GITHUB_TOKEN`
3. Opens a PR into `sample-service-config` bumping `chart/values.yaml` `image.tag` to the new SHA, authenticated via a GitHub App (`APP_ID` + `APP_PRIVATE_KEY` secrets)

The image tag bump is the only thing that triggers the promotion pipeline — all delivery config lives in `sample-service-config`.

## Key conventions

- CI pushes immutable `:<sha>` tags only — `:latest` is never published. The GHCR package must be public so spoke clusters can pull without credentials.
- The GitHub App used by CI must be installed on both `sample-service` and `sample-service-config` with **Contents: read/write** and **Pull requests: read/write**.
- Required repo secrets: `APP_ID` and `APP_PRIVATE_KEY`.
