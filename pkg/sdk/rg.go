package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

func getResourceGroupClient() (resources.GroupsClient, error) {
	groupClient := resources.NewGroupsClient(config.SubID)
	auth, _ := auth.GetResourceManagementAuthorizer()
	groupClient.Authorizer = auth
	//aksClient.AddToUserAgent(config.UserAgent())
	groupClient.PollingDuration = time.Hour * 1
	return groupClient, nil
}

func validateRG(client resources.GroupsClient, name string) (bool, error) {
	r, err := client.CheckExistence(context.Background(), name)
	if err != nil {
		return false, err
	}
	fmt.Printf("%v", r)
	fmt.Printf("%d", r.StatusCode)

	return false, err
}
