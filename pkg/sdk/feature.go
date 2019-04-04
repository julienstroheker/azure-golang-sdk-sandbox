package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-12-01/features"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

type Feature struct {
	client *features.Client
}

func newFeatureClient() (*Feature, error) {
	featuresClient := features.NewClient(config.SubID)
	a, _ := auth.GetResourceManagementAuthorizer()
	featuresClient.Authorizer = a
	featuresClient.PollingDuration = time.Hour * 1
	featClient := &Feature{
		client: &featuresClient,
	}
	return featClient, nil
}

func (feat *Feature) getFeature(featureName, RPNamespace string) (features.Result, error) {
	r, err := feat.client.Get(context.Background(), RPNamespace, featureName)
	if err != nil {
		return features.Result{}, err
	}
	data, _ := json.Marshal(r)
	fmt.Printf("Feature Result: %s", data)
	return r, nil
}
