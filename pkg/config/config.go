package config

import (
	"github.com/Azure/go-autorest/autorest"
)

// Information loaded from the authorization file to identify the client
type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

var (
	cloudEnv            = "AzurePublicCloud"
	StorageAccountName  = "StorageAccountName"
	StorageAccountKey   = "StorageAccountKey"
	ResourceGroupNameSA = "overworry-advisive"
	ResourceName        = "needlewoman-knotlike"
	Location            = "canadacentral"
	BlobName            = "adminkubeconfigju"
	namespace           string
)

var (
	clientData clientInfo
	authorizer autorest.Authorizer
	TenantID   = "TenantID"
	SubID      = "SubID"
	// Julien E2E
	ClientID = "ClientID"
	Secret   = "Secret"
	// juliens2certs100.eastus.cloudapp.azure.com
	ClientID2 = "ClientID2"
	Secret2   = "Secret2"
	TenantID2 = "TenantID2"
	SubID2    = "SubID2"

	ArmAuthorizer      autorest.Authorizer
	KeyvaultAuthorizer autorest.Authorizer
)
