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

type Company struct {
	Id int `json:"id"`
}

type CurrentCollaborator struct {
	Id      int     `json:"id"`
	Company Company `json:"company"`
}

type LoginData struct {
	Firstname           string              `json:"firstname"`
	Lastname            string              `json:"lastname"`
	Token               string              `json:"token"`
	RefreshToken        string              `json:"refresh_token"`
	CurrentCollaborator CurrentCollaborator `json:"current_collaborator_base"`
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

	jsonData = []byte(fmt.Sprintf(`{
    "filters": [
        {
            "_type": "CollaboratorFreedayFilter",
            "name": "date_start",
            "criteria": {
                "_type": "DateRange",
                "min": "2024-10-28T00:00:00Z",
                "max": "2024-12-01T00:00:00Z"
            }
        },
        {
            "_type": "CollaboratorFreedayFilter",
            "name": "company_id",
            "criteria": {
                "_type": "StringValue",
                "value": %d
            }
        },
        {
            "_type": "CollaboratorFreedayFilter",
            "name": "collaborator_ids",
            "criteria": {
                "_type": "ListStringValue",
                "values": [
                    %d
                ]
            }
        }
    ],
    "offset": 0,
    "limit": 25,
    "sort": {
        "field": "date_start",
        "direction": "ASC"
    }
}`, responseObject.Data.CurrentCollaborator.Company.Id, responseObject.Data.CurrentCollaborator.Id))
	req, err = http.NewRequest("POST", "https://rh.silae.fr/api/V1/collaborators/freedays", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+responseObject.Data.Token)

	response, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	body, _ = io.ReadAll(response.Body)

	fmt.Println(string(body))
}
