package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ApiResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type LoginData struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	var jsonData = []byte(`{
        "username": "",
        "password": ""
    }`)
	req, err := http.NewRequest("POST", "https://rh.silae.fr/auth-api/login", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	body, _ := io.ReadAll(response.Body)

	var responseObject ApiResponse[LoginData]
	json.Unmarshal(body, &responseObject)

	fmt.Println(responseObject.Data)
}
