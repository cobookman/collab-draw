package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type IpAddr struct {
	Ip string
}

var (
	hostIp IpAddr
)

func init() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://metadata/computeMetadata/v1/instance"+
		"/network-interfaces/0/access-configs/0/external-ip", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error while talking to metadata server, assuming localhost")
		hostIp.Ip = "localhost"
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	hostIp.Ip = string(body)
}

func HostIpAddr() IpAddr {
	return hostIp
}
