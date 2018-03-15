package main

import (
	"database/sql"
)

type account struct {
	ID       int    `json:"id" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (a *account) getAccount(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM accounts WHERE id=$1", a.ID).Scan(&a.ID, &a.Username, &a.Password)
}

func (a *account) getAccountByUsername(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM accounts WHERE username=$1", a.Username).Scan(&a.ID, &a.Username, &a.Password)
}

func (a *account) insertAccount(db *sql.DB) error {
	hashedPassword, _ := HashPassword(a.Password)
	err := db.QueryRow("INSERT INTO accounts (username, password) VALUES ($1, $2) RETURNING id",
		a.Username,
		hashedPassword).Scan(&a.ID)
	if err != nil {
		return err
	}
	return nil
}
