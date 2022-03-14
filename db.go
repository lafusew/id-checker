package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

type SqlBorrower struct {
	id        int
	firstName sql.NullString
	lastName  sql.NullString
	birthDate sql.NullString
}

type Borrower struct {
	id        int
	firstName string
	lastName  string
	birthDate string
}

func connectDatabase() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PWD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, err
}

func queryIdentity(db *sql.DB) (*sql.Rows, error) {
	fmt.Print("\nSQL QUERY STARTED 🏁 \n\n")

	results, err := db.Query("SELECT id, firstName, lastName, birthDate FROM client WHERE id > 22199")
	if err != nil {
		return nil, err
	}

	return results, err
}

func rowsToStructs(rows *sql.Rows) ([]Borrower, error) {

	var borrowers []Borrower
	for rows.Next() {
		var sqlBorrower SqlBorrower
		var borrower Borrower

		err := rows.Scan(&sqlBorrower.id, &sqlBorrower.firstName, &sqlBorrower.lastName, &sqlBorrower.birthDate)
		if err != nil {
			return borrowers, nil
		}

		borrower.id = sqlBorrower.id

		if sqlBorrower.firstName.Valid {
			borrower.firstName = sqlBorrower.firstName.String
		}

		if sqlBorrower.birthDate.Valid {
			borrower.birthDate = sqlBorrower.birthDate.String
		}

		if sqlBorrower.lastName.Valid {
			borrower.lastName = sqlBorrower.lastName.String
		}

		borrowers = append(borrowers, borrower)
	}

	return borrowers, nil
}
