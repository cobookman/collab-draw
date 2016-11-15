package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ipaddr string
)

func HostIp() string {
	if len(ipaddr) == 0 {
		ipaddr = getHostIp()
	}
	return ipaddr
}

func getHostIp() string {
	url := "http://metadata/computeMetadata/v1/" +
		"/instance/network-interfaces/0/access-configs/0/external-ip"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Metadata-Flavor", "Google")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// assuming localhost
		return "127.0.0.1"
	}

	ip, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		// Error getting body
		return ""
	}
	return string(ip)
}
