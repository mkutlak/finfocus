# Research: Scaleway SKU and Region Resolution

## R1: Which properties name the size of a Scaleway resource

- **Decision**: `type`, then `nodeType`.
- **Rationale**: The pulumiverse Scaleway provider schema (read
  2026-11-03) gives `type` to `instance/server:Server`,
  `kubernetes/cluster:Cluster`, `loadbalancers/loadBalancer:LoadBalancer`,
  `network/publicGateway:PublicGateway`, and `instance/ip:Ip`, and gives
  `nodeType` to `kubernetes/pool:Pool`, `databases/instance:Instance`, and
  `redis/cluster:Cluster`. None of the nine priced types has both.
- **Alternatives considered**: Adding `nodeType` to the generic default
  branch. Rejected because it would change resolution for every provider
  without a case, such as Kubernetes resources.

## R2: Block volumes have no size property

- **Decision**: A fixed SKU `volume` from the well-known SKU table, keyed by
  the type segment `block/volume`.
- **Rationale**: `block/volume:Volume` has only `sizeInGb`, `iops`, and
  `zone`. The table already exists for this exact case
  (`ec2/natgateway`, `eks/cluster`). The plugin picks the storage tier from
  the `iops` tag.
- **Alternatives considered**: Deriving `sbs_5k` or `sbs_15k` from `iops` in
  core. Rejected because it is pricing knowledge (Principle I).

## R3: Region order and environment fallback

- **Decision**: `zone`, then `region`, then `SCW_DEFAULT_ZONE`, then
  `SCW_DEFAULT_REGION`, applied only when the normalized provider is
  `scaleway`.
- **Rationale**: Pools carry both `zone` and `region`, and the zone is the
  narrower location; the catalog prices instances per zone. The provider and
  the Scaleway CLI read these two variables, so a CI job that previews the
  stack already has them. The AWS fallback uses the same scoped pattern
  (`TestResolveSKUAndRegion_AWSRegionFallbackScope`).
- **Alternatives considered**: Reading `scaleway:zone` from stack config.
  Rejected because core does not read provider or stack config for any
  cloud.

## R4: Cache key

- **Decision**: When the properties give no region, the projected cache key's
  region segment takes the provider's environment region
  (`proto.RegionFromEnv`). The key format is unchanged.
- **Rationale**: The key takes its region from the `availabilityZone` or
  `region` property, and the tag digest covers every property, but not an
  environment variable. A Scaleway server with no zone is sent with
  `SCW_DEFAULT_ZONE` as its region, so two runs with different zones would
  otherwise share an entry. The same gap existed for `AWS_REGION`, and this
  closes it too. Entries whose region segment changes miss once.
- **Alternatives considered**: Keying on the fully resolved descriptor
  region. Rejected because resolution logs at debug level on every call, and
  the properties are already in the tag digest.
