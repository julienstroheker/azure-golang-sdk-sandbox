package sdk

import "testing"

func TestValidateRG(t *testing.T) {
	groupClient, _ := getResourceGroupClient()
	validateRG(groupClient, "OS_rg-cypselousp-bgs50udo3h_osa-cobberfell-bgs50udo3h_eastus")
}
