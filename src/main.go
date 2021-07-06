package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

const ENVIRONMENTAL_VARIABLE = "CONFIG"

var _startTime time.Time

type Configuration struct {
	HttpsListenAddress     string
	DomainCertificatePath  string
	PrivateKeyPath         string
	StaticContentDirectory string
}

type HealthInfo struct {
	HostName string
	UpTime   string
}

func loadConfig() (Configuration, error) {
	var configuration Configuration
	var err error
	path := os.Getenv(ENVIRONMENTAL_VARIABLE)
	if len(path) > 0 {
		var data []byte
		data, err = ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal([]byte(data), &configuration)
		}
	} else {
		err_msg := fmt.Sprintf("missing environmental variable '%s'", ENVIRONMENTAL_VARIABLE)
		err = errors.New(err_msg)
	}
	return configuration, err
}

func checkPathError(path string) error {
	_, err := os.Stat(path)
	return err
}

func checkAddrError(addr string) error {
	_, err := regexp.MatchString("\\:\\d+", addr)
	return err
}
func healthCheck(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	upTime := time.Since(_startTime).String()
	healthInfo := HealthInfo{
		hostname,
		upTime,
	}
	json.NewEncoder(w).Encode(healthInfo)
}

func serverFilesOverHttps(add string, key string, cert string, dir string) error {
	log.Print("configurations tls file server")
	fileServer := http.FileServer(http.Dir(dir))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/healthcheck", healthCheck)
	log.Print("starting tls file server")
	err := http.ListenAndServeTLS(add, cert, key, nil)
	log.Print("stopped tls file server")
	return err
}

func main() {
	_startTime = time.Now()

	configuration, err := loadConfig()
	if err != nil {
		log.Fatal("configuration file path: ", err)
		os.Exit(1)
	} else {
		log.Print("configuration was loaded")
	}
	if checkPathError(configuration.PrivateKeyPath) != nil {
		log.Fatal("private key file: ", err)
		os.Exit(2)
	} else {
		log.Print("private key file found")
	}
	if checkPathError(configuration.DomainCertificatePath) != nil {
		log.Fatal("domain certificate file: ", err)
		os.Exit(3)
	} else {
		log.Print("domain certificate file found")
	}
	if checkPathError(configuration.StaticContentDirectory) != nil {
		log.Fatal("static content directory: ", err)
		os.Exit(4)
	} else {
		log.Print("static content directory found")
	}
	if checkAddrError(configuration.HttpsListenAddress) != nil {
		log.Fatal("https listen address: ", err)
		os.Exit(5)
	} else {
		log.Print("https listen address looks good")
	}

	err = serverFilesOverHttps(
		configuration.HttpsListenAddress,
		configuration.PrivateKeyPath,
		configuration.DomainCertificatePath,
		configuration.StaticContentDirectory,
	)
	if err != nil {
		log.Fatal("https endpoint: ", err)
	}
}
