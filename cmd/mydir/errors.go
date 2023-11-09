package main

type UserNotExists struct {
	User string
}

func (e *UserNotExists) Error() string {
	return "The user " + e.User + "do not exist in the system"
}

type UserExists struct {
	User string
}

func (e *UserExists) Error() string {
	return "The user " + e.User + "do not exist in the system"
}

type InvalidPassword struct {
	Message string
}

func (e *InvalidPassword) Error() string {
	return e.Message
}

type TokenExpired struct {
	Token string
}

func (e *TokenExpired) Error() string {
	return "The token " + e.Token + "has expired, login again"
}

type NotOwner struct {
	User  string
	Token string
}

func (e *NotOwner) Error() string {
	return "The token " + e.Token + "is not form user " + e.User
}

type BadAuthHeader struct {
	Message string
}

func (e *BadAuthHeader) Error() string {
	return e.Message
}

type MissingAuthHeader struct {
	Message string
}

func (e *MissingAuthHeader) Error() string {
	return e.Message
}

type MissingDocContent struct {
	Message string
}

func (e *MissingDocContent) Error() string {
	return e.Message
}
