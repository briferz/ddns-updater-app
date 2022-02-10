package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type conf struct {
	queryUrl  string
	apiUrl    string
	hostname  string
	username  string
	password  string
	userAgent string
}

func main() {
	c := getConf()

	log.Print("Retrieving current IP..")
	startTime := time.Now()
	resp, err := http.Get(c.queryUrl)
	elapsedTime := time.Since(startTime)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode >= 300 {
		log.Fatalln("Request to query current IP returned status ", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	currentIp := string(body)

	log.Printf("Retrieved current IP which is %s in %v \n", currentIp, elapsedTime)
	req, err := http.NewRequest(http.MethodGet, c.apiUrl, nil)
	req.SetBasicAuth(c.username, c.password)
	queryValues := req.URL.Query()
	queryValues.Set("hostname", c.hostname)
	queryValues.Set("myip", currentIp)
	req.URL.RawQuery = queryValues.Encode()

	startTime = time.Now()
	resp, err = http.DefaultClient.Do(req)
	elapsedTime = time.Since(startTime)
	if err != nil {
		log.Fatalf("Doing request to update Dynamic IP: %v", err)
	}
	if resp.StatusCode >= 300 {
		log.Fatalln("Request to update Dynamic IP returned status ", resp.Status)
	}
	log.Printf("Updated current IP Address %s successfully at %s in %v\n", currentIp, c.apiUrl, elapsedTime)
}

func retrieveEnvVar(key string, required bool) string {
	v, ok := os.LookupEnv(key)
	if !ok && required {
		log.Fatalf("No environment variable %s was found.", key)
	}
	return v
}

func retrieveEnvVarOrDefault(key string, def string) string {
	v := retrieveEnvVar(key, false)
	if v != "" {
		return v
	}
	return def
}

func getConf() conf {
	defaultUserAgent := fmt.Sprintf("MyProjectCompany %s/%s/%s my.mail@example.com", os.Args[0], runtime.GOOS, runtime.GOARCH)
	return conf{
		queryUrl:  retrieveEnvVar("IP_QUERY_URL", true),
		apiUrl:    retrieveEnvVar("DDNS_API_URL", true),
		hostname:  retrieveEnvVar("DDNS_HOSTNAME", true),
		username:  retrieveEnvVar("DDNS_USERNAME", true),
		password:  retrieveEnvVar("DDNS_PASSWORD", true),
		userAgent: retrieveEnvVarOrDefault("DDNS_USER_AGENT", defaultUserAgent),
	}
}
