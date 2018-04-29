package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//const CERTIFICATE_PATH = "/home/mike/go/src/curl/server.crt"
const CERTIFICATE_PATH = "api.cert.pem"

func isJSON(b []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(b, &js) == nil
}

func prettyPrintJson(b []byte) []byte {
	var pj bytes.Buffer
	err := json.Indent(&pj, b, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json parse error: %v\n", err)
		os.Exit(1)
	}
	return pj.Bytes()
}

func buildClient() *http.Client {

	// load the system cert pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't load system cert pool: %v\n", err)
		os.Exit(1)
	}
	tlsConfig := &tls.Config{RootCAs: certPool}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Load our self signed certificate path
	pemData, err := ioutil.ReadFile(CERTIFICATE_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file open error: %v\n", err)
		os.Exit(1)
	}
	// append our self signed cert to the pool
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	if !ok {
		fmt.Fprintf(os.Stderr, "load cert error: \n")
		os.Exit(1)
	}

	return client
}

func main() {
	url := "https://www.howsmyssl.com/a/check"
	numArgs := len(os.Args[1:])
	if numArgs == 1 {
		url = os.Args[1]
	} else if numArgs > 1 {
		fmt.Println("too many args")
		os.Exit(1)
	}
	client := buildClient()
	rawResp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gurl: %v\n", err)
		os.Exit(1)
	}
	resp, err := ioutil.ReadAll(rawResp.Body)
	rawResp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	if isJSON(resp) {
		fmt.Println("it's json ", rawResp.Status, string(prettyPrintJson(resp)))
	} else {
		fmt.Printf("Response isn't jason  %s %s\n", rawResp.Status, resp)
	}

}
