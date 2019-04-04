package pattern

import "testing"

func TestHttpClient(t *testing.T) {
	httpClient()
	// httpClientWrong()
}

func TestValidateURLisUp(t *testing.T) {
	validateURLisUP("https://openshift2.efb6db3718974eedadf5.dev.azmosa.io")
}
