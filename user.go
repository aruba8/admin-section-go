package main

import "database/sql"

type user struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
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
