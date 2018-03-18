package main

import (
	"database/sql"
)

type user struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
}

func (u *user) getUserById(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM users WHERE id=$1", u.ID).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.MiddleName, &u.Email)
}

func (u *user) insertUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users (firstname, lastname, middlename, email) VALUES ($1, $2, $3, $4) RETURNING id",
		u.FirstName, u.LastName, u.MiddleName, u.Email).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) updateUser(db *sql.DB) error {
	err := db.QueryRow("UPDATE users SET firstname=$1, lastname=$2, middlename=$3, email=$4 WHERE id=$5 RETURNING id, firstname, lastname, middlename, email",
		u.FirstName, u.LastName, u.MiddleName, u.Email, u.ID).Scan(&u.ID, &u.FirstName, &u.LastName, &u.MiddleName, &u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) getUsers(db *sql.DB) ([]user, error) {
	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := []user{}
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.MiddleName, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID);
	return err
}
