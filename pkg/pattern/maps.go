package pattern

import (
	"fmt"
	"regexp"
	"sync"
)

func syncWaitMap() {
	m := make(map[string]error)
	var wg sync.WaitGroup

	certName := []string{"certC", "certB", "certA"}
	//certsOperationsChannel := make(chan map[string]error, len(certName))
	wg.Add(len(certName))
	for _, cert := range certName {
		go fakeAnswer(cert, m, &wg)
	}
	wg.Wait()
	//for range certName {
	//temp := <-certsOperationsChannel
	// for k, v := range m {
	// 	m[k] = v
	// }
	//}
	fmt.Println(m)
}

func syncWaitMapLocked() {
	var certOperationErrorDetails sync.Map
	var wg sync.WaitGroup

	certName := []string{"certC", "certB", "certA"}
	//certsOperationsChannel := make(chan map[string]error, len(certName))
	wg.Add(len(certName))
	for _, cert := range certName {
		go fakeAnswerSync(cert, &certOperationErrorDetails, &wg)
	}
	wg.Wait()
	var certsErrOperations []string

	certOperationErrorDetails.Range(func(certName, certError interface{}) bool {
		fmt.Println(certName)
		fmt.Println(certError)
		certsErrOperations = append(certsErrOperations, fmt.Sprintf("Certificate: %s - Error: %s", certName, certError))

		return true
	})
	fmt.Println(certsErrOperations)
}

func fakeAnswer(cert string, doneWithSuccess map[string]error, wg *sync.WaitGroup) {
	defer wg.Done()
	//c := make(map[string]error)
	//c[cert] = fmt.Errorf("error formated")
	doneWithSuccess[cert] = fmt.Errorf("error formated")
}

func fakeAnswerSync(cert string, doneWithSuccess *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	doneWithSuccess.Store(cert, fmt.Errorf("error formated"))
}

func regex(certURI string) {
	parts := regexp.MustCompile(`^(https://[^/]+)/certificates/([^/]+)(/[^/]+)?$`).FindStringSubmatch(certURI)
	if len(parts) != 4 {
		fmt.Printf("Invalid format URI: %s", certURI)
		return
	}
	//vaultBaseURL := parts[1]
	//certName := parts[2]
	fmt.Printf("parts : %v\n", parts)
	fmt.Printf("len(parts) : %d\n", len(parts))
	fmt.Printf("parts[0] : %v\n", parts[1])
	fmt.Printf("parts[1] : %v\n", parts[1])
	fmt.Printf("parts[2] : %v\n", parts[2])
	fmt.Printf("parts[3] : %v\n", parts[3])
}
