package sdk

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"

	vault "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"

	"github.com/Azure/go-autorest/autorest/azure"
)

const pattern = `^(https://[^/]+)/certificates/([^/]+)(/[^/]+)?$`

var re = regexp.MustCompile(pattern)

type KeyVault struct {
	client      *keyvault.BaseClient
	vaultClient *vault.VaultsClient
	vaultURL    string
}

func newKeyVaultClient(kvName, clientID, clientSecret, tenantID, subID string) (*KeyVault, error) {
	kvClient := keyvault.New()
	vaultClient := vault.NewVaultsClient(subID)
	authkv, _ := auth.GetKeyvaultAuthorizer(clientID, clientSecret, tenantID)
	authvault, _ := auth.GetResourceManagementAuthorizer()
	kvClient.Authorizer = authkv
	vaultClient.Authorizer = authvault
	kvClient.PollingDuration = time.Hour * 1
	vaultClient.PollingDuration = time.Hour * 1
	k := &KeyVault{
		vaultURL:    fmt.Sprintf("https://%s.%s", kvName, azure.PublicCloud.KeyVaultDNSSuffix),
		client:      &kvClient,
		vaultClient: &vaultClient,
	}
	return k, nil
}

func (k *KeyVault) createVault(vaultName, rgName string) error {
	ctx := context.Background()
	tenantIDuuid, err := uuid.FromString(config.TenantID)
	if err != nil {
		log.Fatalf("Error uuid : %v", err)
		return err
	}
	vaultProperties := vault.VaultCreateOrUpdateParameters{
		Location: &config.Location,
		Properties: &vault.VaultProperties{
			TenantID: &tenantIDuuid,
			Sku: &vault.Sku{
				Name:   vault.Standard,
				Family: (func(s string) *string { return &s })("A"),
			},
			AccessPolicies: &[]vault.AccessPolicyEntry{},
		},
	}
	result, err := k.vaultClient.CreateOrUpdate(ctx, rgName, vaultName, vaultProperties)
	if err != nil {
		log.Fatalf("Error creating the vault : %v", err)
		return err
	}
	data, _ := json.Marshal(result)
	log.Printf("Vault created: %s", data)
	return nil
}

func (k *KeyVault) getVault(vaultName, rgName string) error {
	ctx := context.Background()

	result, err := k.vaultClient.Get(ctx, rgName, vaultName)
	if err != nil {
		log.Fatalf("Error getting the vault : %v", err)
		return err
	}
	data, _ := json.Marshal(result)
	log.Printf("Vault : %s", data)
	return nil
}

func (k *KeyVault) deleteVault(vaultName, rgName string) error {
	ctx := context.Background()

	result, err := k.vaultClient.Delete(ctx, rgName, vaultName)
	if err != nil {
		log.Fatalf("Error deleting the vault : %v", err)
		return err
	}
	log.Printf("Vault deleted : %v", result)
	return nil
}

// GetSecret retrieves a secret from keyvault
func (k *KeyVault) getSecretAndParsePem(keyName string) (keyvault.SecretBundle, error) {
	ctx := context.Background()

	secretBundle, err := k.client.GetSecret(ctx, k.vaultURL, keyName, "")
	if err != nil {
		return keyvault.SecretBundle{}, err
	}
	// data, _ := json.Marshal(secretBundle)
	// fmt.Printf("Get Secret parsed : %s\n", data)
	// fmt.Printf("Get Secret value : %s\n", *secretBundle.Value)

	var cert2 *x509.Certificate
	var key2 *rsa.PrivateKey
	data := []byte(*secretBundle.Value)
	fmt.Printf("Get Secret parsed : %s\n", data)
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		switch block.Type {
		case "CERTIFICATE":
			cert2, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				//return err
			}
			fmt.Printf("cert2 : %+v\n", cert2)

		case "PRIVATE KEY":
			k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				//return err
			}
			key2 = k.(*rsa.PrivateKey)
			fmt.Printf("cert2 : %+v\n", key2)
		}
	}
	return secretBundle, nil
}

// GetSecret retrieves a secret from keyvault
func (k *KeyVault) getSecret(keyName string) (keyvault.SecretBundle, error) {
	ctx := context.Background()

	secretBundle, err := k.client.GetSecret(ctx, k.vaultURL, keyName, "")
	if err != nil {
		return keyvault.SecretBundle{}, err
	}
	data, _ := json.Marshal(secretBundle)
	fmt.Printf("Get Secret parsed : %s\n", data)
	fmt.Printf("Get Secret value : %s\n", *secretBundle.Value)
	return secretBundle, nil
}

// GetSecret retrieves a secret from keyvault
func (k *KeyVault) getCertificate(certName string) (keyvault.CertificateBundle, error) {
	ctx := context.Background()

	certBundle, err := k.client.GetCertificate(ctx, k.vaultURL, certName, "")
	if err != nil {
		return keyvault.CertificateBundle{}, err
	}

	return certBundle, nil
}

// GetSecret retrieves a secret from keyvault
func (k *KeyVault) backupAndRestore(certName, restoreVaultURI string) {
	bkCerts, err := k.client.BackupCertificate(context.Background(), k.vaultURL, certName)
	if err != nil {
		log.Fatal("Error Backup")
	}
	restoreCert := keyvault.CertificateRestoreParameters{
		CertificateBundleBackup: bkCerts.Value,
	}
	result, err := k.client.RestoreCertificate(context.Background(), restoreVaultURI, restoreCert)
	if err != nil {
		log.Fatalf("Error Backup : %v", err)
	}
	log.Printf("restore result : %v", result)
}

func (k *KeyVault) importCert(value *string, vaultRestoreURI, certName string) {
	ctx := context.Background()
	certToImport := keyvault.CertificateImportParameters{
		Base64EncodedCertificate: value,
	}
	result, err := k.client.ImportCertificate(ctx, vaultRestoreURI, certName, certToImport)
	if err != nil {
		log.Fatalf("Error Import : %v", err)
	}
	log.Printf("restore result : %v", result)
}

func (k *KeyVault) generateCert(certName string) error {
	ctx := context.Background()
	tags := map[string]*string{
		"SubscriptionID":    (func(s string) *string { return &s })("SubscriptionID"),
		"ResourceGroupName": (func(s string) *string { return &s })("ResourceGroupName"),
		"ResourceName":      (func(s string) *string { return &s })("ResourceName"),
	}
	parameters := keyvault.CertificateCreateParameters{
		CertificatePolicy: &keyvault.CertificatePolicy{
			IssuerParameters: &keyvault.IssuerParameters{
				Name: (func(s string) *string { return &s })("digicert01"),
			},
			KeyProperties: &keyvault.KeyProperties{
				Exportable: (func(s bool) *bool { return &s })(true),
				KeySize:    (func(s int32) *int32 { return &s })(2048),
				KeyType:    keyvault.RSA,
				ReuseKey:   (func(s bool) *bool { return &s })(false),
			},
			SecretProperties: &keyvault.SecretProperties{
				ContentType: (func(s string) *string { return &s })("application/x-pkcs12"),
			},
			X509CertificateProperties: &keyvault.X509CertificateProperties{
				Subject: (func(s string) *string { return &s })("CN=" + certName + ".test.julien.cloudapp.azure.com"),
			},
		},
		Tags: tags,
	}
	result, err := k.client.CreateCertificate(ctx, "tedsfsfdsfdsf", certName, parameters)
	if err != nil {
		log.Fatalf("Error Cert Creation : %v", err)
		return err
	}
	data, _ := json.Marshal(result)
	log.Printf("Cert Creation : %s", data)
	return nil
}

func (k *KeyVault) deleteCert(certId string) error {
	ctx := context.Background()
	parts := re.FindStringSubmatch(certId)
	fmt.Printf("len parts : %d", len(parts))
	fmt.Printf("parts : %s", parts)
	if len(parts) != 4 {
		log.Fatalf("Error Cert Parts : %d", len(parts))
	}
	vaultBaseURL := parts[1]
	certName := parts[2]
	//certStatus := strings.TrimPrefix(parts[3], "/")
	result, err := k.client.DeleteCertificate(ctx, vaultBaseURL, certName)
	if err != nil {
		log.Fatalf("Error Cert Deletion : %v", err)
		return err
	}
	log.Printf("Cert Deletion not parsed : %+v", result)
	data, _ := json.Marshal(result)
	log.Printf("Cert Deletion parsed : %s", data)
	return nil
}

// GetSecret retrieves a secret from keyvault
func (k *KeyVault) getCertificateOperation(certName string) (keyvault.CertificateOperation, error) {
	ctx := context.Background()

	certOperation, err := k.client.GetCertificateOperation(ctx, k.vaultURL, certName)
	if err != nil {
		return keyvault.CertificateOperation{}, err
	}
	data, _ := json.Marshal(certOperation)
	log.Printf("Cert Creation : %s", data)
	return certOperation, nil
}

// Merge Access Policy
func (k *KeyVault) mergeKVAccessPolicyPermissions(old, new vault.Permissions) vault.Permissions {
	blankPermissions := vault.Permissions{}
	var certificateMerged []vault.CertificatePermissions
	var secretMerged []vault.SecretPermissions
	var keyMerged []vault.KeyPermissions
	if new.Certificates != nil {
		for _, certPerm := range *new.Certificates {
			certificateMerged = append(certificateMerged, certPerm)
		}
	}
	if old.Certificates != nil {
		for _, certPerm := range *old.Certificates {
			certificateMerged = append(certificateMerged, certPerm)
		}
	}
	if new.Keys != nil {
		for _, keyPerm := range *new.Keys {
			keyMerged = append(keyMerged, keyPerm)
		}
	}
	if old.Keys != nil {
		for _, keyPerm := range *old.Keys {
			keyMerged = append(keyMerged, keyPerm)
		}
	}
	if new.Secrets != nil {
		for _, secretPerm := range *new.Secrets {
			secretMerged = append(secretMerged, secretPerm)
		}
	}
	if old.Secrets != nil {
		for _, secretPerm := range *old.Secrets {
			secretMerged = append(secretMerged, secretPerm)
		}
	}
	blankPermissions.Certificates = &certificateMerged
	blankPermissions.Secrets = &secretMerged
	blankPermissions.Keys = &keyMerged
	return blankPermissions
}

func (k *KeyVault) getGetKVAccessPolicyTemplate() *vault.Permissions {
	return &vault.Permissions{
		Certificates: &[]vault.CertificatePermissions{
			vault.Get,
		},
		Keys: &[]vault.KeyPermissions{
			vault.KeyPermissionsGet,
		},
	}
}

func (k *KeyVault) getImportKVAccessPolicyTemplate() *vault.Permissions {
	return &vault.Permissions{
		Certificates: &[]vault.CertificatePermissions{
			vault.Import,
		},
		Secrets: &[]vault.SecretPermissions{
			vault.SecretPermissionsGet,
		},
	}
}

func (k *KeyVault) mergeKVAccessPolicies() ([]vault.AccessPolicyEntry, error) {
	blankAccessPolicy := []vault.AccessPolicyEntry{}
	tenantIDuuid, _ := uuid.FromString(config.TenantID)
	old := vault.AccessPolicyEntry{
		ObjectID: (func(s string) *string { return &s })(config.TenantID),
		TenantID: &tenantIDuuid,
		Permissions: &vault.Permissions{
			Certificates: &[]vault.CertificatePermissions{
				vault.Get,
				vault.Import,
			},
			Secrets: &[]vault.SecretPermissions{
				vault.SecretPermissionsGet,
			},
			Keys: &[]vault.KeyPermissions{
				vault.KeyPermissionsGet,
			},
		},
	}
	new := vault.AccessPolicyEntry{
		ObjectID: (func(s string) *string { return &s })(config.TenantID),
		TenantID: &tenantIDuuid,
		Permissions: &vault.Permissions{
			Certificates: &[]vault.CertificatePermissions{
				vault.Import,
			},
		},
	}
	blankAccessPolicy = append(blankAccessPolicy, new)
	blankAccessPolicy = append(blankAccessPolicy, old)
	return blankAccessPolicy, nil
}
