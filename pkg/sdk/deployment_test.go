package sdk

import "testing"

func TestGetDeploy(t *testing.T) {
	dc := getDeploymentClient()
	getDeployment(dc, "55d30728-a303-48a3-a274-b1781dd03479")
}
