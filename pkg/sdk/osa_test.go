package sdk

import (
	"log"
	"testing"
)

func TestOSAList(t *testing.T) {
	osaCLient, err := newOSAClient()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	err = osaCLient.ListOSAinSub()
	if err != nil {
		log.Fatal(err)
	}
}
