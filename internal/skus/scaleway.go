package skus

// scalewayWellKnownSKUMap maps Pulumi module/resource strings to fixed SKU
// labels for Scaleway types that carry no size property. The pricing plugin
// reads the remaining dimensions (for a volume: iops and sizeInGb) from tags.
//
//nolint:gochecknoglobals // Package-level lookup table, initialized once.
var scalewayWellKnownSKUMap = map[string]string{
	"block/volume": "volume",
}

// resolveScalewaySKU returns the fixed SKU label for a Scaleway Pulumi type
// token such as "scaleway:block/volume:Volume", or "" when none applies.
func resolveScalewaySKU(resourceType string) string {
	return scalewayWellKnownSKUMap[extractPulumiSegment(resourceType)]
}
