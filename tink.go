package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getClientToken() {
	data := url.Values{}
	data.Set("client_id", os.Getenv("TINK_CLIENT"))
	data.Set("client_secret", os.Getenv("TINK_SECRET"))
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "authorization:grant")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.tink.com/api/v1/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(string(body))
}