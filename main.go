package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type IdentityCheck struct {
	mbClientId int
	mbFullName string
	mbBirthDate string
	tinkFullName string
	tinkBirthDate string
	levenshteinDistance int
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

	token, err := getClientToken()
	if err != nil {
		panic(err.Error())
	}

	err = initLocalDatabase()
	if err != nil {
		panic(err.Error())
	}

	for i, v := range borrowers {
		code, err := getUserCode(token, v.id)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if code == "" {
			continue 
		}
	
		userToken, err := getUserTokenFromCode(code)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	
		availableIdentities, err := getAvailableUserIdentity(userToken)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		err = insertLocalDatabase(availableIdentities, v.id)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Fprint(os.Stdout, "\r \r")
		fmt.Printf("%d / %d", i, len(borrowers))
	}
}

// func compare(bor Borrower, tinkIdentity Identity) string {
// 	mbFullName := fmt.Sprintf("%s %s", bor.firstName, bor.lastName)
// 	tinkFullName := tinkIdentity.Name

// 	return fmt.Sprintf("\nmb  : %s\ntink: %s \n", mbFullName, tinkFullName)
// }


// func createAndOpenTxt() *os.File {
// 	err := ioutil.WriteFile("compare.txt", []byte("Quick print of all data:\n"), 0644)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	file, err := os.OpenFile("compare.txt", os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return file
// }