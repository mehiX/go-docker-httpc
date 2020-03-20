package client

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// DockerHTTP Sends a request to the docker daemon and returns the response
func DockerHTTP(q string, method string) (string, error) {

	log.Printf("Docker -> query: %s, method: %s", q, method)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	req, err := http.NewRequest(method, "http://unix"+q, nil)
	if nil != err {
		return "", err
	}

	resp, err := client.Do(req)
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	return string(body), nil
}
