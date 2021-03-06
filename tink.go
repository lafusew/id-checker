package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Token struct {
	Token string `json:"access_token"`
	Type string `json:"token_type"`
	Hint string `json:"id_hint"`
	ExpiresIn float64 `json:"expires_in"`
	Refresh string `json:"refresh_token"`
	Scope string `json:"scope"`
}

type UserCode struct {
	Code string `json:"code"`
}

type Identity struct {
	Name string `json:"name"`
	Ssn string `json:"ssn"`
	BirthDate string `json:"dateOfBirth"`
	BankProvider string `json:"providerName"`
}

type IdentitiesResp struct {
	Identities []Identity `json:"availableIdentityData"`
}

var db *sql.DB

func getClientToken() (Token, error){
	data := url.Values{}
	data.Set("client_id", os.Getenv("TINK_CLIENT"))
	data.Set("client_secret", os.Getenv("TINK_SECRET"))
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "authorization:grant")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.tink.com/api/v1/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return Token{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()

	var b Token

	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return Token{}, err
	}

	return b, nil
}

func getUserCode(token Token, userId int) (string, error){
	data := url.Values{}
	data.Set("external_user_id", fmt.Sprint(userId))
	data.Set("scope", "accounts:read,transactions:read,identity:read")
	
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.tink.com/api/v1/oauth/authorization-grant", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.Type, token.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var b UserCode

	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return "", err
	}

	return b.Code, err
}

func getUserTokenFromCode(code string) (Token, error){
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("TINK_CLIENT"))
	data.Set("client_secret", os.Getenv("TINK_SECRET"))
	data.Set("grant_type", "authorization_code")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.tink.com/api/v1/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return Token{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(req)
	if err != nil {
		return Token{}, err
	}

	defer res.Body.Close()

	var b Token

	err = json.NewDecoder(res.Body).Decode(&b);
	if err != nil {
		return Token{}, err
	}

	return b, err
}

func getAvailableUserIdentity(token Token) (IdentitiesResp, error){
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.tink.com/api/v1/identities", strings.NewReader(""))
	if err != nil {
		return IdentitiesResp{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.Type, token.Token))

	res, err := client.Do(req)
	if err != nil {
		 return IdentitiesResp{}, err
	}
	defer res.Body.Close()

	var b IdentitiesResp
	
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return IdentitiesResp{}, err
 	}

	return b, err
}

func initLocalDatabase() error {
	var err error

	db, err = sql.Open("sqlite3", "./tink.db")

	if err != nil {
		return err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tink (
		"id" INTEGER NOT NULL,
		"name" TEXT,
		"ssn" TEXT,
		"dateOfBirth" TEXT,
		"providerName" TEXT
	);`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		return err
	}

	statement.Exec()

	return err
}

func insertLocalDatabase(tinkRes IdentitiesResp, mbId int) error {
	insertCredSQL := `INSERT INTO tink (id, name, ssn, dateOfBirth, providerName) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertCredSQL)
	if err != nil {
		return err
	} else {
		defer statement.Close()
	}

	for _, v := range tinkRes.Identities {
		if err != nil {
			return err
		}

		_, err := statement.Exec(mbId, v.Name, v.Ssn, v.BirthDate, v.BankProvider)
		if err != nil {
			return err
		}
	}


	return nil
}