package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type conf struct {
	queryUrl string
	apiUrl   string
	hostname string
	username string
	password string
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

func retrieveEnvVar(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("No environment variable %s was found.", key)
	}
	return v
}

func getConf() conf {
	return conf{
		queryUrl: retrieveEnvVar("IP_QUERY_URL"),
		apiUrl:   retrieveEnvVar("DDNS_API_URL"),
		hostname: retrieveEnvVar("DDNS_HOSTNAME"),
		username: retrieveEnvVar("DDNS_USERNAME"),
		password: retrieveEnvVar("DDNS_PASSWORD"),
	}
}
