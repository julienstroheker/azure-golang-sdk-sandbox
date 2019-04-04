package sdk

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
)

func TestCreateAppCertsService(t *testing.T) {
	ctx := context.Background()
	certClient, err := getCertsClient()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	myAppCert := web.AppServiceCertificateOrder{
		Location: (func(s string) *string { return &s })("global"),
		Name:     (func(s string) *string { return &s })("testviasdk"),
		AppServiceCertificateOrderProperties: &web.AppServiceCertificateOrderProperties{
			ProductType:       web.CertificateProductType(web.StandardDomainValidatedSsl),
			DistinguishedName: (func(s string) *string { return &s })("CN=sdk.julien.work"),
		},
	}
	result, err := certClient.CreateOrUpdate(ctx, "testviasdk", "testviasdk", myAppCert)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	result.WaitForCompletionRef(ctx, certClient.BaseClient.Client)
	response, err := result.Future.GetResult(certClient)
	_, err = result.Result(certClient)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	fmt.Printf("Result : %v \n", result.Future)
	fmt.Printf("Result : %v \n", response)

	mycert, err := certClient.Get(ctx, "testviasdk", "testviasdk")
	fmt.Printf("Token : %s \n", *mycert.AppServiceCertificateOrderProperties.DomainVerificationToken)
}
