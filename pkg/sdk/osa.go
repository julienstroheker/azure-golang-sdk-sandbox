package sdk

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2018-09-30-preview/containerservice"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

type OSA struct {
	client *containerservice.OpenShiftManagedClustersClient
}

func newOSAClient() (*OSA, error) {
	osaClient := containerservice.NewOpenShiftManagedClustersClient(config.SubID)
	authosa, _ := auth.GetResourceManagementAuthorizer()
	osaClient.Authorizer = authosa
	osaClient.PollingDuration = time.Hour * 1
	osa := &OSA{
		client: &osaClient,
	}
	return osa, nil
}

func (osa *OSA) ListOSAinSub() error {

	i, err := osa.client.List(context.Background())
	if err != nil {
		log.Fatal("Error: " + err.Error())
		return err
	}
	// for i.NotDone() {
	// 	r = append(r, i.Values())
	// 	i.Next()
	// }
	d := i.Values()
	data, _ := json.Marshal(d)
	log.Printf("ListOSA: %s", data)
	return nil
}

func (osa *OSA) ListOSAinRG(rgName string) error {
	var r []containerservice.OpenShiftManagedCluster
	i, err := osa.client.ListByResourceGroupComplete(context.Background(), rgName)
	if err != nil {
		log.Fatal("Error: " + err.Error())
		return err
	}
	for i.NotDone() {
		r = append(r, i.Value())
		i.Next()
	}
	data, _ := json.Marshal(r)
	log.Printf("ListOSA in RG: %s", data)
	return nil
}

func renderOSACluster(name string) containerservice.OpenShiftManagedCluster {
	fqdn := (func(s string) *string { return &s })(name + ".westcentralus.cloudapp.azure.com")
	return containerservice.OpenShiftManagedCluster{
		Location: (func(s string) *string { return &s })("westcentralus"),
		Name:     (func(s string) *string { return &s })(name),
		OpenShiftManagedClusterProperties: &containerservice.OpenShiftManagedClusterProperties{
			OpenShiftVersion: (func(s string) *string { return &s })("v3.11"),
			Fqdn:             fqdn,
			AgentPoolProfiles: &[]containerservice.OpenShiftManagedClusterAgentPoolProfile{
				containerservice.OpenShiftManagedClusterAgentPoolProfile{
					Name:       (func(s string) *string { return &s })("compute"),
					Count:      (func(i int32) *int32 { return &i })(5),
					VMSize:     containerservice.OpenShiftContainerServiceVMSize(containerservice.StandardD4sV3),
					SubnetCidr: (func(s string) *string { return &s })("10.0.0.0/24"),
					Role:       containerservice.OpenShiftAgentPoolProfileRole(containerservice.Compute),
					OsType:     containerservice.OSType(containerservice.Linux),
				},
				containerservice.OpenShiftManagedClusterAgentPoolProfile{
					Name:       (func(s string) *string { return &s })("infra"),
					Count:      (func(i int32) *int32 { return &i })(2),
					VMSize:     containerservice.OpenShiftContainerServiceVMSize(containerservice.StandardD4sV3),
					SubnetCidr: (func(s string) *string { return &s })("10.0.0.0/24"),
					Role:       containerservice.OpenShiftAgentPoolProfileRole(containerservice.Infra),
					OsType:     containerservice.OSType(containerservice.Linux),
				},
			},
			MasterPoolProfile: &containerservice.OpenShiftManagedClusterMasterPoolProfile{
				Count:      (func(i int32) *int32 { return &i })(3),
				VMSize:     containerservice.OpenShiftContainerServiceVMSize(containerservice.StandardD4sV3),
				SubnetCidr: (func(s string) *string { return &s })("10.0.0.0/24"),
			},
			NetworkProfile: &containerservice.NetworkProfile{
				VnetCidr: (func(s string) *string { return &s })("10.0.0.0/8"),
			},
			RouterProfiles: &[]containerservice.OpenShiftRouterProfile{
				containerservice.OpenShiftRouterProfile{
					Name: (func(s string) *string { return &s })("default"),
				},
			},
			AuthProfile: &containerservice.OpenShiftManagedClusterAuthProfile{
				IdentityProviders: &[]containerservice.OpenShiftManagedClusterIdentityProvider{
					containerservice.OpenShiftManagedClusterIdentityProvider{
						Name: (func(s string) *string { return &s })("Azure AD"),
						Provider: containerservice.OpenShiftManagedClusterAADIdentityProvider{
							ClientID: (func(s string) *string { return &s })(config.ClientID),
							Secret:   (func(s string) *string { return &s })(config.Secret),
							TenantID: (func(s string) *string { return &s })(config.TenantID),
							Kind:     containerservice.KindAADIdentityProvider,
						},
					},
				},
			},
		},
	}
}
