package ms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	clientID = "9eb2aae0-c030-49a8-866b-304064770509"
	tenantID = "9d21d004-0c5d-4069-bf84-4c799d627d43"
)

func GetAccessToken() (string, error) {
	scope := "offline_access User.Read Calendars.ReadWrite"
	deviceCodeURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode?mkt=fr-FR", tenantID)
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	// 1. Demande du code d’appareil
	resp, err := http.PostForm(deviceCodeURL, map[string][]string{
		"client_id": {clientID},
		"scope":     {scope},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var deviceCodeResp DeviceCodeResp
	if err := json.Unmarshal(body, &deviceCodeResp); err != nil {
		return "", err
	}

	fmt.Println(deviceCodeResp.Message)

	// 2. Polling pour le token d’accès
	for {
		time.Sleep(time.Duration(deviceCodeResp.Interval) * time.Second)

		data := url.Values{}
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		data.Set("client_id", clientID)
		data.Set("device_code", deviceCodeResp.DeviceCode)

		req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusOK {
			var tokenResp TokenResponse
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				return "", err
			}

			return tokenResp.AccessToken, nil
		}

		fmt.Println("En attente de l’authentification de l’utilisateur... Veuillez continuer sur votre navigateur.")
	}
}

func CreateOutlookEvent(accessToken string) error {
	eventUrl := "https://graph.microsoft.com/v1.0/me/events"
	event := map[string]interface{}{
		"subject": "Absence",
		"start": map[string]string{
			"dateTime": time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
			"timeZone": "Europe/Paris",
		},
		"end": map[string]string{
			"dateTime": time.Now().Add(25 * time.Hour).Format("2006-01-02T15:04:05"),
			"timeZone": "Europe/Paris",
		},
		"showAs": "oof",
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", eventUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("bad API HTTP status: " + resp.Status)
	}

	fmt.Println("Événement créé avec succès")
	return nil
}
