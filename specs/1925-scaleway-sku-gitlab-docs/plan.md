# Implementation Plan: Scaleway SKU and Region Resolution, GitLab CI Docs

**Branch**: `1925-scaleway-sku-gitlab-docs` | **Date**: 2026-11-03 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/1925-scaleway-sku-gitlab-docs/spec.md`

## Summary

Add a `scaleway` case to the shared SKU and region resolver so that Scaleway
pools, databases, Redis clusters, and block volumes pass strict pre-flight
validation and reach a Scaleway pricing plugin. Add one Scaleway row to the
well-known SKU table (`block/volume` → `volume`). Rewrite the GitLab CI
example and fix two stale Docker guide passages.

## Technical Context

**Language/Version**: Go (version from `go.mod`)
**Primary Dependencies**: finfocus-spec `pluginsdk/mapping` (`ExtractSKU`, `ExtractRegion` with explicit keys); no new dependency
**Storage**: N/A
**Testing**: testify, table tests in `internal/proto/adapter_test.go` and `internal/skus/resolve_test.go`
**Target Platform**: Linux, macOS, Windows (pure Go, environment reads only)
**Project Type**: Single Go module (CLI)
**Performance Goals**: N/A (two map lookups per resource)
**Constraints**: No new field sent to plugins; cache key unchanged; no other provider's resolution changes
**Scale/Scope**: About 30 lines of production code, 2 docs pages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Plugin-First Architecture**: Orchestration only. Core learns Scaleway
  property names (`type`, `nodeType`, `zone`, `region`) and one fixed SKU
  label, the same kind of knowledge as the existing `aws`, `azure`, and
  `gcp` cases and the AWS well-known SKU table. Prices stay in the plugin.
- [x] **Test-Driven Development**: Table tests are written first and fail
  before the `scaleway` case exists.
- [x] **Cross-Platform Compatibility**: Environment variable reads only.
- [x] **Documentation Integrity**: The CI/CD integration and Docker guides
  are updated in this feature. No other page describes Scaleway.
- [x] **Protocol Stability**: No proto or SDK change.
- [x] **Implementation Completeness**: No stubs or TODOs.
- [x] **Persistence Model**: No persistent state.
- [x] **Web Parity**: The spec records Web UI Fit as Not applicable with a
  reason (no new command, flag, or view).
- [x] **Quality Gates**: `make lint`, `make test`, `make docs-lint`.
- [x] **Multi-Repo Coordination**: The plugin lives in
  `finfocus-plugin-scaleway-public`. It depends on this change for four of
  its nine types and needs no spec bump.

**Violations Requiring Justification**: None.

## Project Structure

### Documentation (this feature)

```text
specs/1925-scaleway-sku-gitlab-docs/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── resolution.md
└── tasks.md             # /speckit.tasks
```

### Source Code (repository root)

```text
internal/proto/
├── adapter.go           # resolveSKUAndRegion: add case "scaleway" and its env fallback
└── adapter_test.go      # Scaleway resolution table + env-scope test
internal/skus/
├── resolve.go           # ResolveSKU: add case "scaleway"
├── scaleway.go          # new: scalewayWellKnownSKUMap {"block/volume": "volume"}
└── resolve_test.go      # Scaleway rows
docs/src/content/docs/deployment/
├── cicd-integration.md  # GitLab section rewritten
└── docker.md            # stale "planned" and manual-plugin text fixed
```

**Structure Decision**: Single module. The change stays in the two
packages that already own SKU and region resolution.

## Complexity Tracking

No violations.
