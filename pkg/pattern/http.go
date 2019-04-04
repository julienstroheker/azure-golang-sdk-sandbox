package pattern

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

func httpClient() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	response, err := http.Get("https://openshift.efb6db3718974eedadf5.dev.azmosa.io")
	if err != nil {
		fmt.Println(err)
	}
	data, _ := json.Marshal(response.TLS)
	fmt.Printf("response %s", data)
	// fmt.Printf("response %+v", response)

}

func httpClientWrong() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	response, err := http.Get("https://julientest1000.eastus.cloudapp.azure.com/")
	if err != nil {
		fmt.Println(err)
		data, _ := json.Marshal(err)
		fmt.Printf("error %s", data)
		return
	}
	data, _ := json.Marshal(response.TLS)
	fmt.Printf("response %s", data)
	// fmt.Printf("response %+v", response)

}

func validateURLisUP(uri string) {
	response, err := http.Get(uri)
	if err != nil {
		fmt.Printf("err (%s)", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		fmt.Println("All good 200")
	} else {
		fmt.Printf("Not Good : %d", response.StatusCode)
	}
	// data, _ := json.Marshal(response)
	// fmt.Printf("http response (%s)", data)
}
