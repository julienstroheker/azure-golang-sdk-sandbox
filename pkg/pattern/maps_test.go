package pattern

import "testing"

func TestSyncWaitMap(t *testing.T) {
	syncWaitMap()
}

func TestSyncWaitMapLocked(t *testing.T) {
	syncWaitMapLocked()
}

func TestRegex(t *testing.T) {
	regex("https://myvault.vault.azure.net/certificates/foo")
}
