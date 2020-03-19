package service

import (
	"DockerHttpClient/data"
	"errors"
)

type user struct {
	username string
	password string
	profile  data.User
}

var users []user

func init() {
	users = make([]user, 0)
	users = append(users, user{
		username: "mihai",
		password: "password",
		profile: data.User{
			Username: "mihai",
			Email:    "mihai@devops.org",
			Role:     "admin",
		},
	})

	users = append(users, user{
		username: "user1",
		password: "userpass",
		profile: data.User{
			Username: "user1",
			Email:    "user1@devops.org",
			Role:     "user",
		},
	})
}

// Login checks if the credentials are valid
func Login(username string, passwd string) (*data.User, error) {

	for _, u := range users {
		if username == u.username && passwd == u.password {
			return &u.profile, nil
		}
	}

	return nil, errors.New("Unknown user credentials")
}
