package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var secret = RandomString(12)

func validateAccount(username string, password string, db *sql.DB) bool {
	acc := account{Username: username}
	if err := acc.getAccountByUsername(db); err != nil {
		return false
	}
	if !CheckPasswordHash(password, acc.Password) {
		return false
	}
	return true
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14);
	return string(bytes), err
}

func RandomString(strlen int) string {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, strlen)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

type ClaimsStruct struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	jwt.StandardClaims
}

type tokenStruct struct {
	Token string `json:"token"`
}
