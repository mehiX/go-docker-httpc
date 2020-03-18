package client

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
)

// DockerHttp Sends a request to the docker daemon and returns the response
func DockerHttp(q string, method string) (string, error) {
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
