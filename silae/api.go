package silae

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserData(username, password string) (*UserData, error) {
	payload, err := json.Marshal(credentials{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://rh.silae.fr/auth-api/login", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad API HTTP status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse[UserData]
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Status != "ok" {
		return nil, errors.New("can not log to Silae API")
	}

	return &apiResp.Data, nil
}

func GetFreedays(userData *UserData) (*FreedaysData, error) {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	nextMonthDate := currentDate.AddDate(0, 1, 0)
	payload := RequestPayload{
		Filters: []Filter{
			{
				Type: "CollaboratorFreedayFilter",
				Name: "date_start",
				Criteria: DateRangeCriteria{
					Type: "DateRange",
					Min:  currentDate.Format(time.RFC3339),
					Max:  nextMonthDate.Format(time.RFC3339),
				},
			},
			{
				Type: "CollaboratorFreedayFilter",
				Name: "company_id",
				Criteria: StringValueCriteria{
					Type:  "StringValue",
					Value: userData.CurrentCollaborator.Company.Id,
				},
			},
			{
				Type: "CollaboratorFreedayFilter",
				Name: "collaborator_ids",
				Criteria: ListStringValueCriteria{
					Type:   "ListStringValue",
					Values: []int{userData.CurrentCollaborator.Id},
				},
			},
		},
		Offset: 0,
		Limit:  25,
		Sort: Sort{
			Field:     "date_start",
			Direction: "ASC",
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://rh.silae.fr/api/V1/collaborators/freedays", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userData.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad API HTTP status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse[FreedaysData]
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Status != "ok" {
		return nil, errors.New("can not get freedays from Silae API")
	}

	return &apiResp.Data, nil
}