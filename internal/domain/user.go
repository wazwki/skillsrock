package domain

import "errors"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var UserNotFound = errors.New("User not found")
