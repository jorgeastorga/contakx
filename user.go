package main

import "golang.org/x/crypto/bcrypt"

type User struct {
  ID string
  Username string
  Email string
  HashedPassword string
}


const (
    passwordLength = 8
    hashCost = 10
    userIDLength
)

func NewUser(username, email, password string) (User, error){
    user := User{
        Email:    email,
        Username: username,
    }
    if username == "" {
        return user, errNoUsername
    }

    if email == "" {
        return user, errNoEmail
    }

    if password == "" {
        return user, errNoPassword
    }

    if len(password) < passwordLength {
        return user, errPasswordTooShort
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte (password), hashCost)

    user.HashedPassword = string(hashedPassword)
    user.ID = GenerateID("usr", userIDLength)

    return user, err
}
