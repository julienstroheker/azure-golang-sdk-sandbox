package sdk

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/julienstroheker/azure-golang-sdk-sandbox/pkg/config"
)

func BlobStorage() {

	if len(config.StorageAccountName) == 0 || len(config.StorageAccountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}
	credential, err := azblob.NewSharedKeyCredential(config.StorageAccountName, config.StorageAccountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerName := config.Location + "-" + config.ResourceGroupNameSA + "-" + config.ResourceName

	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", config.StorageAccountName, containerName, config.BlobName))

	blobURL := azblob.NewBlobURL(*u, p)
	ctx := context.Background() // This example uses a never-expiring context

	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)

	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	fmt.Printf("%s", downloadedData.String())

}
