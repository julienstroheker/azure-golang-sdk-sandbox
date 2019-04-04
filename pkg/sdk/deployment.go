package sdk

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

func getDeploymentClient() resources.DeploymentsClient {
	deploymentsc := resources.NewDeploymentsClient(config.SubID)
	auth, _ := auth.GetResourceManagementAuthorizer()
	deploymentsc.Authorizer = auth
	return deploymentsc
}

func getDeployment(client resources.DeploymentsClient, rg string) {
	var one int32 = 3
	i, _ := client.ListByResourceGroupComplete(context.Background(), rg, "", &one)
	for i.NotDone() {
		v := i.Value()
		fmt.Printf("Name %s\n", *v.Name)
		fmt.Printf("Provision State %s\n", *v.Properties.ProvisioningState)
		fmt.Printf("-----")
		i.Next()
	}

}
