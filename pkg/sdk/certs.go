package sdk

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/julienstroheker/sandbox/pkg/auth"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

func getCertsClient() (web.AppServiceCertificateOrdersClient, error) {
	certClient := web.NewAppServiceCertificateOrdersClient(config.SubID)
	auth, _ := auth.GetResourceManagementAuthorizer()
	certClient.Authorizer = auth
	//aksClient.AddToUserAgent(config.UserAgent())
	certClient.PollingDuration = time.Hour * 1
	return certClient, nil
}
