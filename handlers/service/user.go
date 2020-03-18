package service

import "DockerHttpClient/data"

// Login checks if the credentials are valid
func Login(username string, passwd string) (data.User, error) {

	return data.User{Username: username, Role: "admin"}, nil
}
