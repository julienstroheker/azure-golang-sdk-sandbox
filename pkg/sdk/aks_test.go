package sdk

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestAKSGet(t *testing.T) {
	ctx := context.Background()
	aksCLient, err := getAKSClient()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	mc, err := aksCLient.Get(ctx, "gdfgfdgfdgfdgdfgfd", "gdfgfdgfdgfdgdfgfd")
	fmt.Printf("Cluster FQDN : %s", *mc.Fqdn)
	mcLisRG, err := aksCLient.ListByResourceGroup(ctx, "tgdfgfdgfdgfdgdfgfd")
	fmt.Printf("%v", mcLisRG)
	lista := listAKS(&aksCLient)
	fmt.Printf("%v", lista)
	fmt.Println()

	fmt.Printf("Result : %v \n", mc.Status)
}
