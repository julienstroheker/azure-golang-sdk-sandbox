package sdk

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
	"github.com/julienstroheker/sandbox/pkg/utils"
)

const kvName = "osa-int-certs"
const kvRGName = "rp-common-int"

func TestGetSecret(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	bundle, err := keyClient.getSecret("gdfgfdgfdgfdgdfgfd")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Secret : %v", bundle)
}

func TestGetSecretAndParsePem(t *testing.T) {
	keyClient, err := newKeyVaultClient("kv9ea872c6b1f845aaaaf3", config.ClientID2, config.Secret2, config.TenantID2, config.SubID2)
	if err != nil {
		log.Fatal(err)
	}
	bundle, err := keyClient.getSecretAndParsePem("dasdsadasda503c96feb124f1d8aaadevdasdsadsado")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Secret : %v", bundle)
}

func TestGetCert(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	certBundle, err := keyClient.getCertificate("dasdsadsafafadfdfdfds")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Certificate : %v", certBundle)
}

func TestBackupRestore(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	bundle, err := keyClient.getSecret("test-manual-julien")
	if err != nil {
		log.Fatal(err)
	}

	pem, err := utils.ExtractPEMFromKeyvaultSecretBundle(&bundle)

	log.Printf("Retrieved secret '%s' from keyvault", *bundle.Value)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.X509KeyPair(pem, pem)
	if err != nil {
		log.Printf("Cannot load X509 key pair: %v", err)
	}

	keyClient.backupAndRestore("test-manual-julien", "https://julientest.vault.azure.net")
}

func TestGetAndImport(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	keyClientImport, err := newKeyVaultClient("kv98917c3e1ef04b35899c", config.ClientID2, config.Secret2, config.TenantID2, config.SubID2)
	if err != nil {
		log.Fatal(err)
	}
	bundle, err := keyClient.getSecret("dasdffsdgfsgfsgfdgfdgfdgd")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Secret : %v", bundle.Value)
	certBundle, err := keyClientImport.getCertificate("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("certBundle : %v", certBundle)
	keyClientImport.importCert(bundle.Value, "https://kv98917c3e1ef04b35899c.vault.azure.net", "dasdffsdgfsgfsgfdgfdgfdgd")
}

func TestCreateCert(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	err = keyClient.generateCert("CreatedWithGO2")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeleteCert(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	err = keyClient.deleteCert("https://dasdsadasdasdasdsadsad.vault.azure.net/certificates/test-julien-delete2")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCertOperation(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	op, err := keyClient.getCertificateOperation("testestestesttesosaio")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Operation : %s", *op.Status)
}

func TestCreateKV(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	err = keyClient.createVault("8Vault-Created-WithGO", "julienCertsTests")
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetKV(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	err = keyClient.getVault(kvName, kvRGName)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMergeKVPolicyPerm(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	old := keyClient.getGetKVAccessPolicyTemplate()
	new := keyClient.getImportKVAccessPolicyTemplate()
	mergedKVPolicy := keyClient.mergeKVAccessPolicyPermissions(*new, *old)
	data, _ := json.Marshal(mergedKVPolicy)
	fmt.Printf("%s", data)
}

func TestMergeKVAccessPolicies(t *testing.T) {
	keyClient, err := newKeyVaultClient(kvName, config.ClientID, config.Secret, config.TenantID, config.SubID)
	if err != nil {
		log.Fatal(err)
	}
	toto, err := keyClient.mergeKVAccessPolicies()
	data, _ := json.Marshal(toto)
	fmt.Printf("%s", data)
}
