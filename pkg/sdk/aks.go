package sdk

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2018-09-30-preview/containerservice"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

func getAKSClient() (containerservice.ManagedClustersClient, error) {
	aksClient := containerservice.NewManagedClustersClient(config.SubID)
	auth, _ := auth.GetResourceManagementAuthorizer()
	aksClient.Authorizer = auth
	//aksClient.AddToUserAgent(config.UserAgent())
	aksClient.PollingDuration = time.Hour * 1
	return aksClient, nil
}

func listAKS(client *containerservice.ManagedClustersClient) []containerservice.ManagedCluster {
	//var r []containerservice.ManagedCluster
	i, err := client.List(context.Background())
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	// for i.NotDone() {
	// 	r = append(r, i.Response())
	// 	i.Next()
	// }
	return i.Values()
}
