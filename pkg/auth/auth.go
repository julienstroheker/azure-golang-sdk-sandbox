package auth

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// OAuthGrantType specifies which grant type to use.
type OAuthGrantType int

const (
	OAuthGrantTypeServicePrincipal OAuthGrantType = iota
)

// GrantType returns what grant type has been configured.
func grantType() OAuthGrantType {
	return OAuthGrantTypeServicePrincipal
}

func GetAuthorizerForResource(grantType OAuthGrantType, resource string) (autorest.Authorizer, error) {

	var a autorest.Authorizer
	var err error

	switch grantType {

	case OAuthGrantTypeServicePrincipal:
		azureEnv, _ := azure.EnvironmentFromName("AzurePublicCloud")
		authorizationServerURL := azureEnv.ActiveDirectoryEndpoint
		oauthConfig, err := adal.NewOAuthConfig(
			authorizationServerURL, config.TenantID)
		if err != nil {
			return nil, err
		}

		token, err := adal.NewServicePrincipalToken(
			*oauthConfig, config.ClientID, config.Secret, resource)
		if err != nil {
			return nil, err
		}
		a = autorest.NewBearerAuthorizer(token)

	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	return a, err
}

func GetResourceManagementAuthorizer() (autorest.Authorizer, error) {
	if config.ArmAuthorizer != nil {
		return config.ArmAuthorizer, nil
	}

	var a autorest.Authorizer
	var err error

	a, err = GetAuthorizerForResource(
		OAuthGrantTypeServicePrincipal, "https://management.azure.com/")

	if err == nil {
		// cache
		config.ArmAuthorizer = a
	} else {
		// clear cache
		config.ArmAuthorizer = nil
	}
	return config.ArmAuthorizer, err
}

// GetKeyvaultAuthorizer gets an OAuthTokenAuthorizer for use with Key Vault
// keys and secrets. Note that Key Vault *Vaults* are managed by Azure Resource
// Manager.
func GetKeyvaultAuthorizer(clientID, clientSecret, tenantID string) (autorest.Authorizer, error) {
	if config.KeyvaultAuthorizer != nil {
		return config.KeyvaultAuthorizer, nil
	}

	// BUG: default value for KeyVaultEndpoint is wrong
	azureEnv, _ := azure.EnvironmentFromName("AzurePublicCloud")
	vaultEndpoint := strings.TrimSuffix(azureEnv.KeyVaultEndpoint, "/")
	// BUG: alternateEndpoint replaces other endpoints in the configs below
	alternateEndpoint, _ := url.Parse(
		"https://login.windows.net/" + tenantID + "/oauth2/token")

	var a autorest.Authorizer
	var err error

	switch grantType() {
	case OAuthGrantTypeServicePrincipal:
		oauthconfig, err := adal.NewOAuthConfig(
			azureEnv.ActiveDirectoryEndpoint, tenantID)
		if err != nil {
			return a, err
		}
		oauthconfig.AuthorizeEndpoint = *alternateEndpoint

		token, err := adal.NewServicePrincipalToken(
			*oauthconfig, clientID, clientSecret, vaultEndpoint)
		if err != nil {
			return a, err
		}

		a = autorest.NewBearerAuthorizer(token)
	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	if err == nil {
		config.KeyvaultAuthorizer = a
	} else {
		config.KeyvaultAuthorizer = nil
	}

	return config.KeyvaultAuthorizer, err
}
