package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type SqlBorrower struct {
	id        int
	firstName sql.NullString
	lastName  sql.NullString
}

type Borrower struct {
	id        int
	firstName string
	lastName  string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	
	db, err := connectDatabase()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := queryIdentity(db)
	if err != nil {
		panic(err.Error())
	}

	borrowers, err := rowsToStructs(rows)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(borrowers)
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

	results, err := db.Query("SELECT id, firstName, lastName FROM client")
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

		err := rows.Scan(&sqlBorrower.id, &sqlBorrower.firstName, &sqlBorrower.lastName)
		if err != nil {
			return borrowers, nil
		}

		borrower.id = sqlBorrower.id

		if !sqlBorrower.firstName.Valid {
			borrower.firstName = ""
		} else {
			borrower.firstName = sqlBorrower.firstName.String
		}

		if !sqlBorrower.lastName.Valid {
			borrower.lastName = ""
		} else {
			borrower.lastName = sqlBorrower.lastName.String
		}

		borrowers = append(borrowers, borrower)
	}

	return borrowers, nil
}