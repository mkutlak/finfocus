# Feature Specification: Scaleway SKU and Region Resolution, GitLab CI Docs

**Feature Branch**: `1925-scaleway-sku-gitlab-docs`
**Created**: 2026-11-03
**Status**: Implemented
**Input**: User description: "Scaleway SKU and region resolution in core, plus GitLab CI docs, so a Scaleway pricing plugin (the separate finfocus-plugin-scaleway-public repo) receives every priced Scaleway resource, and a self-hosted GitLab pipeline can run finfocus from the published image."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Every priced Scaleway resource reaches the plugin (Priority: P1)

An operator has a Pulumi stack on Scaleway and installs a Scaleway pricing
plugin. They run `finfocus cost projected` on a preview. Kubernetes pools,
managed databases, Redis clusters, and block volumes are sent to the plugin
like servers and load balancers, and are not stopped by pre-flight validation
with a `$0` `VALIDATION:` row.

**Why this priority**: Today those four types carry no property that core
reads as a SKU, so validation rejects them before any plugin is asked. A
Scaleway estimate then leaves out its largest line items (node pools and
databases).

**Independent Test**: Build the projected descriptor for each Scaleway
resource type from a preview fixture and check that its SKU and region are
non-empty and that strict validation passes.

**Acceptance Scenarios**:

1. **Given** a `scaleway:kubernetes/pool:Pool` with `nodeType: PRO2-S` and
   `zone: fr-par-1`, **When** core builds its projected request, **Then** the
   SKU is `PRO2-S` and the region is `fr-par-1`.
2. **Given** a `scaleway:databases/instance:Instance` or
   `scaleway:redis/cluster:Cluster` with a `nodeType`, **When** core builds
   the request, **Then** the SKU is that `nodeType`.
3. **Given** a `scaleway:instance/server:Server` with `type: DEV1-S`, **When**
   core builds the request, **Then** the SKU is `DEV1-S`, as it is today.
4. **Given** a `scaleway:block/volume:Volume`, which has no type property,
   **When** core builds the request, **Then** the SKU is `volume`, and the
   plugin still receives `sizeInGb` and `iops` in the tags.
5. **Given** a resource with both `zone` and `region`, **When** core builds
   the request, **Then** the region field holds the zone.

---

### User Story 2 - Creates that take their zone from provider config (Priority: P2)

A program sets the zone once on the Scaleway provider (or through the
`SCW_DEFAULT_ZONE` variable), not on each resource. A preview of a new
resource then has no `zone` or `region` input. The CI job sets
`SCW_DEFAULT_ZONE` (or `SCW_DEFAULT_REGION`), and the resource is priced in
that zone.

**Why this priority**: It only affects resources being created whose
location is not set on the resource, but without it those rows are rejected
for an empty region.

**Independent Test**: Resolve a Scaleway resource with no location
properties while `SCW_DEFAULT_ZONE` is set, then with only
`SCW_DEFAULT_REGION` set, then with neither.

**Acceptance Scenarios**:

1. **Given** no `zone` or `region` property and `SCW_DEFAULT_ZONE=nl-ams-1`,
   **When** core resolves the region, **Then** the region is `nl-ams-1`.
2. **Given** no location property, no `SCW_DEFAULT_ZONE`, and
   `SCW_DEFAULT_REGION=pl-waw`, **When** core resolves the region, **Then**
   the region is `pl-waw`.
3. **Given** an AWS or Azure resource with no region and `SCW_DEFAULT_ZONE`
   set, **When** core resolves the region, **Then** the Scaleway variables
   are ignored.

---

### User Story 3 - Run finfocus in self-hosted GitLab CI (Priority: P2)

A platform engineer on self-hosted GitLab wants a merge request pipeline
that prices the Pulumi preview. They copy the documented two-job example.
One job runs `pulumi preview --json` in their Pulumi image and keeps
`plan.json`. A second job runs on the published finfocus image, installs a
plugin, and prices the plan.

**Why this priority**: The current GitLab example installs a release that is
many versions old through a pinned script, and the Docker guide still says
the image is planned and that plugins must be copied in by hand.

**Independent Test**: Run the documented second job's commands in the
published image against a fixture plan and get a priced table. Lint the docs.

**Acceptance Scenarios**:

1. **Given** the CI/CD integration guide, **When** a reader opens its GitLab
   section, **Then** it shows the two-job example with the image, an empty
   entrypoint, `finfocus plugin install`, and `finfocus cost projected
   --pulumi-json plan.json`, and no `v0.3.1` install script.
2. **Given** the Docker guide, **When** a reader looks for how to add a
   plugin, **Then** it shows `finfocus plugin install` and does not say that
   building the image is planned.

---

### Edge Cases

- A Scaleway resource whose `type` or `nodeType` is the Pulumi unknown value
  at preview time: the existing sentinel filter drops that property, so the
  SKU is empty and validation rejects the row, as it does for other
  providers.
- A resource with both `type` and `nodeType`: `type` wins.
- `SCW_DEFAULT_ZONE` set to an empty string: treated as unset.
- A resource that has no location property and no Scaleway variable set
  still falls back to a region inherited from a referenced resource, as
  every provider does today.
- A `scaleway:` type that has no SKU property and is not a block volume
  (for example a private network) keeps an empty SKU and is rejected by
  validation as it is today. The plugin cannot price it anyway.
- Legacy `scaleway:index/*` tokens: the type segment does not match
  `block/volume`, so those volumes are not mapped.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: For a resource whose provider is `scaleway`, core MUST take the
  SKU from the `type` property, then from `nodeType`.
- **FR-002**: When neither property gives a SKU, core MUST use the existing
  well-known SKU table. For Scaleway that table MUST map the type segment
  `block/volume` to `volume` and nothing else.
- **FR-003**: For a `scaleway` resource, core MUST take the region from the
  `zone` property, then `region`, then the `SCW_DEFAULT_ZONE` environment
  variable, then `SCW_DEFAULT_REGION`.
- **FR-004**: The Scaleway environment variables MUST NOT affect any other
  provider, and AWS, Azure, GCP, and generic resolution MUST be unchanged.
- **FR-005**: This MUST apply wherever core already resolves SKU and region
  (projected, recommendations, actual-cost tags, and cross-resource
  reference identity), because those paths share one resolver.
- **FR-006**: Core MUST NOT send any new field to plugins. The projected
  cache key format MUST NOT change, but its region segment MUST reflect a
  region taken from the environment, so runs with different
  `SCW_DEFAULT_ZONE` (or `AWS_REGION`) values do not share an entry.
- **FR-007**: Core MUST NOT add Scaleway prices, product names, or unit
  conversions. Pricing stays in the plugin.
- **FR-008**: The GitLab section of the CI/CD integration guide MUST show a
  two-job example: a preview job that saves `plan.json` as an artifact, and
  a job on `ghcr.io/rshade/finfocus:<version>` with an empty entrypoint that
  installs a plugin and prices the plan. The stale install-script example
  MUST be removed.
- **FR-009**: The Docker guide MUST describe `finfocus plugin install` for
  adding plugins and MUST NOT describe the image or its Dockerfile as
  planned.

### Key Entities

- **Projected descriptor**: The request core sends to a plugin for one
  resource: provider, type, SKU, region, and the resource's properties as
  tags. This feature changes only how SKU and region are filled for
  Scaleway.
- **Well-known SKU table**: The core table of fixed SKUs for resource types
  that have no SKU property. It gains one Scaleway row.

## Web UI Fit *(mandatory)*

- **Disposition**: Not applicable. No command, flag, or view is added.
  Scaleway rows are priced through the same resolver and appear on the
  existing Estimate, Overview, and Cost pages exactly as CLI rows do.
- **Web data**: None new. The pages already show whatever the plugin prices.
- **Layout fit**: None. No view changes.
- **Additions**: None.
- **Render evidence**: Not applicable. There is no new or changed page.
- **Follow-up issue scope**: No `web-ui` issue. There is no UX remainder.

## Assumptions

- The Pulumi Scaleway provider (pulumiverse) names its types
  `scaleway:<module>/<resource>:<Class>`. The location properties are `zone`
  and `region`, and the size properties are `type` or `nodeType`.
- The Scaleway CLI and provider read `SCW_DEFAULT_ZONE` and
  `SCW_DEFAULT_REGION`, so a CI job that previews a stack already sets them.
- A zone (for example `fr-par-1`) is a valid value for the descriptor's
  region field. The plugin derives the region from the zone when it needs
  one.
- Following the existing per-cloud property cases for AWS, Azure, and GCP
  meets Principle I, because core learns property names and not prices.
- Terraform state `scaleway_*` resources already map to provider `scaleway`
  with camelCase keys, so they take the same resolution (`node_type` becomes
  the SKU). Terraform support is otherwise out of scope.

## Out of Scope

- The Scaleway pricing plugin itself (a separate repository).
- Actual costs from the Scaleway Billing API.
- Terraform `scaleway_*` types and legacy `scaleway:index/*` tokens.
- A `registry.json` entry for the Scaleway plugin.
- A GitLab CI/CD component or merge request comments (a separate GitLab
  project).
- Adding tools such as `pulumi` or `jq` to the published image.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For a preview fixture with one resource of each of the nine
  priced Scaleway types (server, pool, cluster, block volume, load balancer,
  database, Redis, public gateway, flexible IP), 9 of 9 pass strict
  pre-flight validation (up from 5 of 9).
- **SC-002**: The existing AWS, Azure, and GCP SKU and region resolution tests
  pass unchanged.
- **SC-003**: A Scaleway resource with no location property and
  `SCW_DEFAULT_ZONE` set is priced in that zone, and AWS or Azure resources
  ignore the variable in the same run.
- **SC-004**: A reader can copy the GitLab example into a pipeline and get a
  priced table for a preview without installing anything except a plugin.
- **SC-005**: The docs pass `make docs-lint`, and the GitLab example no
  longer references `v0.3.1`.
