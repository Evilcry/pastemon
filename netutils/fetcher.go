package netutils

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Get func
// wrap ProcessHTTP
func Get(url string) []byte {
	return ProcessHTTP(http.Get(url))
}

// ProcessHTTP func
// process HTTP response
func ProcessHTTP(resp *http.Response, err error) []byte {
	if err != nil {
		log.Println("[-] Could not access url.")
		return []byte("")
	}

	if resp.StatusCode != 200 {
		log.Printf("[-] Received HTTP error %d.\n", resp.StatusCode)
		return []byte("")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[-] Could not read response")
		return []byte("")
	}

	return body
}
