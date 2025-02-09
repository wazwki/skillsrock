package domain

import "errors"

var UserNotFound = errors.New("User not found")

type User struct {
	ID       int
	Name     string
	Password string
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func UserToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func UserRequestToUser(user *UserRequest) *User {
	return &User{
		Name:     user.Name,
		Password: user.Password,
	}
}
