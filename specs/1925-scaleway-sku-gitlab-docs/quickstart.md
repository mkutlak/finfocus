# Quickstart: Scaleway SKU and Region Resolution

```bash
go test -run 'TestResolveSKUAndRegion_Scaleway|TestPrepareProjectedDescriptor_Scaleway' ./internal/proto/
go test -run TestProjectedCacheKeyIncludesEnvRegion ./internal/engine/
go test ./internal/skus/
make lint && make test && make docs-lint
```

Manual check against a Scaleway preview. Pre-flight validation runs only
when a plugin is installed, so install the recorder plugin into a temporary
home and read the requests it records:

```bash
make build build-recorder
home=$(mktemp -d) out=$(mktemp -d)
mkdir -p "$home/plugins/recorder/0.1.0"
cp bin/finfocus-plugin-recorder plugins/recorder/plugin.manifest.json "$home/plugins/recorder/0.1.0/"
FINFOCUS_HOME="$home" FINFOCUS_RECORDER_OUTPUT_DIR="$out" \
  ./bin/finfocus cost projected --pulumi-json plan.json
grep -h '"sku"\|"region"' "$out"/*Supports*.json
```

Every Scaleway resource appears with a non-empty SKU and region.
