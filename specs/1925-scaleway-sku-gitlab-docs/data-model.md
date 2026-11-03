# Data Model: Scaleway SKU and Region Resolution

No stored data changes. The only shape involved is the existing projected
descriptor (`pbc.ResourceDescriptor`).

| Field | Scaleway source after this feature |
| --- | --- |
| `provider` | Unchanged (`scaleway`) |
| `resource_type` | Unchanged |
| `sku` | `type`, else `nodeType`, else the well-known SKU table (`block/volume` → `volume`), else empty |
| `region` | `zone`, else `region`, else `SCW_DEFAULT_ZONE`, else `SCW_DEFAULT_REGION`, else the existing referenced-resource inheritance |
| `tags` | Unchanged: every property, including `sizeInGb`, `iops`, `size`, `clusterSize`, `isHaCluster` |

Validation is unchanged: strict validation still rejects an empty SKU or
region.
