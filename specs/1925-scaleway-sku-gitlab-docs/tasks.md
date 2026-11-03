# Tasks: Scaleway SKU and Region Resolution, GitLab CI Docs

**Input**: Design documents from `specs/1925-scaleway-sku-gitlab-docs/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/resolution.md

**Tests**: Required. Constitution Principle II (TDD): each test task comes
before its implementation and must fail first.

## Phase 1: Setup

None. No new module, dependency, or tool.

## Phase 2: Foundational

None. Both user stories touch only existing packages.

## Phase 3: User Story 1 - Every priced Scaleway resource reaches the plugin (P1) 🎯 MVP

**Goal**: Scaleway pools, databases, Redis clusters, and block volumes get a
non-empty SKU and region and pass strict validation.

**Independent Test**: `go test -run 'TestResolveSKUAndRegion_Scaleway|TestResolveSKU' ./internal/proto/ ./internal/skus/`

- [X] T001 [P] [US1] Add Scaleway rows to `TestResolveSKU` in `internal/skus/resolve_test.go`: `scaleway:block/volume:Volume` → `volume`, `scaleway:instance/server:Server` → empty, a legacy `scaleway:index/blockVolume:BlockVolume` → empty. Run it and confirm the volume row fails.
- [X] T002 [P] [US1] Add `TestResolveSKUAndRegion_ScalewayProperties` (parallel table) in `internal/proto/adapter_test.go` covering the property rows of `contracts/resolution.md`: server `type`+`zone`, pool `nodeType`+`zone`+`region` (zone wins), database `nodeType`+`region`, Redis `nodeType`+`zone`, block volume with no type → `volume`, and `type` winning over `nodeType`. Also assert that `ValidateProjectedDescriptor(PrepareProjectedDescriptor(...))` returns nil for each of the nine priced types from spec SC-001. Run it and confirm the `nodeType` and volume rows fail.
- [X] T003 [US1] Create `internal/skus/scaleway.go` with `scalewayWellKnownSKUMap` (`"block/volume": "volume"`) and `resolveScalewaySKU`, which reuses `extractPulumiSegment`. Add `case "scaleway"` to `ResolveSKU` in `internal/skus/resolve.go` and update its doc comment ("currently aws and scaleway").
- [X] T004 [US1] Add `case "scaleway"` to `resolveSKUAndRegion` in `internal/proto/adapter.go`: `sku = mapping.ExtractSKU(properties, "type", "nodeType")`, falling back to `skus.ResolveSKU(provider, resourceType, properties)`, and `region = mapping.ExtractRegion(properties, "zone", "region")`. Use a `scalewayProvider` constant beside `awsProvider`. T001 and T002 pass.

**Checkpoint**: US1 is complete; a Scaleway plugin receives all nine types.

## Phase 4: User Story 2 - Creates that take their zone from provider config (P2)

**Goal**: A Scaleway resource with no location property takes
`SCW_DEFAULT_ZONE`, then `SCW_DEFAULT_REGION`. No other provider reads them.

**Independent Test**: `go test -run TestResolveSKUAndRegion_ScalewayRegionFallbackScope ./internal/proto/`

- [X] T005 [US2] Add `TestResolveSKUAndRegion_ScalewayRegionFallbackScope` in `internal/proto/adapter_test.go`, modeled on `TestResolveSKUAndRegion_AWSRegionFallbackScope` (not parallel; it uses `t.Setenv`, so give it a plain comment, not `//nolint:paralleltest`). Rows: zone env used; region env used when the zone env is unset; an explicit `zone` beats the env; an empty `SCW_DEFAULT_ZONE` is treated as unset; AWS and Azure resources ignore both variables. Run it and confirm it fails.
- [X] T006 [US2] In `internal/proto/adapter.go`, after the per-cloud switch, add the Scaleway-only fallback beside the AWS block: when `region == "" && cloud == scalewayProvider`, use `os.Getenv("SCW_DEFAULT_ZONE")`, then `SCW_DEFAULT_REGION`. Extend the `resolveSKUAndRegion` doc comment with the Scaleway order. T005 passes.

**Checkpoint**: US1 and US2 pass together.

## Phase 5: User Story 3 - Run finfocus in self-hosted GitLab CI (P2)

**Goal**: A GitLab reader copies a working two-job example, and the Docker
guide matches the published image.

**Independent Test**: `make docs-lint`. Then run the documented second job's
commands in `ghcr.io/rshade/finfocus:latest` against
`test/fixtures/plans/` and get a table.

- [X] T007 [P] [US3] Rewrite the `## GitLab CI` section of `docs/src/content/docs/deployment/cicd-integration.md`. Job `pulumi-preview` runs `pulumi preview --json > plan.json` in the user's Pulumi image (stage `preview`, `artifacts: paths: [plan.json]`, MR pipelines). Job `finfocus-cost` uses `image: {name: ghcr.io/rshade/finfocus:<version>, entrypoint: [""]}`, `needs: [pulumi-preview]`, then runs `finfocus plugin install <plugin>` and `finfocus cost projected --pulumi-json plan.json`. Name the variables the preview needs (`PULUMI_ACCESS_TOKEN`, cloud credentials) and say that a plugin is installed per run. Remove the `v0.3.1` install-script snippet.
- [X] T008 [P] [US3] Fix `docs/src/content/docs/deployment/docker.md`: remove the notes saying no Dockerfile exists and that container support is planned (top note and the local-build note), drop "(planned - requires Dockerfile)", and replace `### Plugin Installation` with `finfocus plugin install` run in a container whose plugin directory is mounted. Check the build-variable list against `docker/Dockerfile`, which declares no `ARG`; correct or remove it. Change the `v0.3.0` pull example to `v<version>`.

**Checkpoint**: All stories are done.

## Phase 6: Polish & Cross-Cutting

Web Parity phase omitted: spec Web UI Fit is Not applicable.

- [X] T009 Run `go vet ./...`, `make lint`, `make test`, and `make docs-lint`, and fix any failures.
- [X] T010 Build (`make build build-recorder`), install the recorder plugin into `FINFOCUS_HOME=$(mktemp -d)`, and run `./bin/finfocus cost projected --pulumi-json <scaleway fixture>` with a nine-type Scaleway preview fixture. Confirm that the recorder received every resource with a non-empty SKU and region (SC-001); without a plugin, validation does not run.
- [X] T011 Review follow-up: in `generateProjectedCostResourceKey` (`internal/engine/engine.go`), fill an empty region segment from `proto.RegionFromEnv` so an environment region is part of the projected cache key (FR-006), with `TestProjectedCacheKeyIncludesEnvRegion` in `internal/engine/engine_cache_test.go`. Fix the GitLab preview image entrypoint and the Docker plugin-volume example, and add the Scaleway sources to `docs/src/content/docs/troubleshooting.md`.

## Dependencies & Execution Order

- US1 (T001-T004) first. T003 depends on T001, and T004 depends on T002 and T003.
- US2 (T005-T006) depends on T004, because both edit `resolveSKUAndRegion`.
- US3 (T007-T008) is independent of US1 and US2 and can run in parallel with them.
- Polish (T009-T010) runs after all stories.

## Parallel Example

```text
T001 (internal/skus/resolve_test.go)  ||  T002 (internal/proto/adapter_test.go)  ||  T007 (cicd-integration.md)  ||  T008 (docker.md)
```

## Implementation Strategy

MVP is US1: it unblocks four of the nine plugin types. US2 then covers creates
whose location comes from provider config. US3 is docs and ships with them.
