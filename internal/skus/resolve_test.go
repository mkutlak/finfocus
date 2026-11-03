package skus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rshade/finfocus/internal/skus"
)

func TestResolveSKU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		provider     string
		resourceType string
		properties   map[string]string
		expected     string
	}{
		{
			name:         "AWS EKS cluster resolves to cluster SKU",
			provider:     "aws",
			resourceType: "aws:eks/cluster:Cluster",
			expected:     "cluster",
		},
		{
			name:         "AWS EKS addon resolves to addon SKU",
			provider:     "aws",
			resourceType: "aws:eks/addon:Addon",
			expected:     "addon",
		},
		{
			name:         "unknown AWS resource returns empty",
			provider:     "aws",
			resourceType: "aws:s3/bucket:Bucket",
			expected:     "",
		},
		{
			name:         "unknown provider returns empty",
			provider:     "azure",
			resourceType: "azure:compute/virtualMachine:VirtualMachine",
			expected:     "",
		},
		{
			name:         "empty provider returns empty",
			provider:     "",
			resourceType: "aws:eks/cluster:Cluster",
			expected:     "",
		},
		{
			name:         "empty resource type returns empty",
			provider:     "aws",
			resourceType: "",
			expected:     "",
		},
		{
			name:         "case-insensitive provider matching",
			provider:     "AWS",
			resourceType: "aws:eks/cluster:Cluster",
			expected:     "cluster",
		},
		{
			name:         "Scaleway block volume resolves to volume SKU",
			provider:     "scaleway",
			resourceType: "scaleway:block/volume:Volume",
			properties:   map[string]string{"sizeInGb": "50", "iops": "5000"},
			expected:     "volume",
		},
		{
			name:         "Scaleway server has no well-known SKU",
			provider:     "scaleway",
			resourceType: "scaleway:instance/server:Server",
			expected:     "",
		},
		{
			name:         "legacy Scaleway index token is not mapped",
			provider:     "scaleway",
			resourceType: "scaleway:index/blockVolume:BlockVolume",
			expected:     "",
		},
		{
			name:         "properties are passed through but not used for well-known types",
			provider:     "aws",
			resourceType: "aws:eks/cluster:Cluster",
			properties:   map[string]string{"name": "my-cluster"},
			expected:     "cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := skus.ResolveSKU(tt.provider, tt.resourceType, tt.properties)
			assert.Equal(t, tt.expected, result)
		})
	}
}
