package utils

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"

	"golang.org/x/crypto/pkcs12"
)

// ExtractPEMFromKeyvaultSecretBundle converts a keyvault secret bundle to PEM
func ExtractPEMFromKeyvaultSecretBundle(secret *keyvault.SecretBundle) ([]byte, error) {
	if secret.ContentType == nil {
		return nil, fmt.Errorf("ContentType cannot be nil")
	}
	if secret.Value == nil {
		return nil, fmt.Errorf("Value cannot be nil")
	}
	switch *secret.ContentType {
	case "application/x-pem-file":
		return []byte(*secret.Value), nil
	case "application/x-pkcs12":
		return p12ToPEM(*secret.Value)
	default:
		return nil, fmt.Errorf("Unknown content type %s", *secret.ContentType)
	}
}

func p12ToPEM(base64P12 string) ([]byte, error) {
	p12, err := base64.StdEncoding.DecodeString(base64P12)
	if err != nil {
		return nil, err
	}
	// pfx in keyvault does not have password
	blocks, err := pkcs12.ToPEM(p12, "")
	if err != nil {
		return nil, err
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	return pemData, nil
}
