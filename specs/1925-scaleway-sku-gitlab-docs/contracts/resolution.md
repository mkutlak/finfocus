# Contract: Scaleway SKU and Region Resolution

Entry point: `proto.ResolveSKUAndRegion(ctx, provider, resourceType, properties)`
(and the unexported `resolveSKUAndRegion` it wraps). It is shared by
`PrepareProjectedDescriptor`, recommendations target resources,
actual-cost tag enrichment, and `engine` reference identity.

Applies when `resourcetype.NormalizeProvider(provider) == "scaleway"`.

| Input properties | Environment | SKU | Region |
| --- | --- | --- | --- |
| `type=DEV1-S`, `zone=fr-par-1` | none | `DEV1-S` | `fr-par-1` |
| `nodeType=PRO2-S`, `size=3`, `zone=fr-par-2`, `region=fr-par` | none | `PRO2-S` | `fr-par-2` |
| `nodeType=DB-DEV-S`, `region=fr-par` | none | `DB-DEV-S` | `fr-par` |
| `sizeInGb=50`, `iops=5000`, `zone=nl-ams-1` (type `scaleway:block/volume:Volume`) | none | `volume` | `nl-ams-1` |
| `type=DEV1-S` | `SCW_DEFAULT_ZONE=pl-waw-1` | `DEV1-S` | `pl-waw-1` |
| `type=DEV1-S` | `SCW_DEFAULT_REGION=pl-waw` only | `DEV1-S` | `pl-waw` |
| `type=DEV1-S`, `zone=fr-par-1` | `SCW_DEFAULT_ZONE=pl-waw-1` | `DEV1-S` | `fr-par-1` |

For any other provider, `SCW_DEFAULT_ZONE` and `SCW_DEFAULT_REGION` are not
read.

`skus.ResolveSKU("scaleway", resourceType, props)` returns `volume` for the
type segment `block/volume` and empty for every other segment.

`proto.RegionFromEnv(provider)` returns that same environment region (or
""). `generateProjectedCostResourceKey` uses it for the region segment when
the `availabilityZone` and `region` properties are empty.
