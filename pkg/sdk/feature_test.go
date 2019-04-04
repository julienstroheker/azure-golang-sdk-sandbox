package sdk

import (
	"fmt"
	"testing"
)

func TestGetFeatures(t *testing.T) {
	featureClient, _ := newFeatureClient()
	_, err := featureClient.getFeature("Microsoft.ContainerService/SaveOSATestConfig", "Microsoft.ContainerService")
	if err != nil {
		fmt.Printf("Error fetching feature : %v", err)
	}
}
