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

	"github.com/spf13/viper"
)

const (
	clientID        = "9eb2aae0-c030-49a8-866b-304064770509"
	tenantID        = "9d21d004-0c5d-4069-bf84-4c799d627d43"
	EventDateLayout = "2006-01-02T15:04:05"
)

func GetAccessToken() (string, error) {
	// 1. Utilisation du refresh_token si possible
	refreshToken := viper.GetString("ms.refresh_token")
	if refreshToken != "" {
		token, err := obtainAccessToken(refreshToken, true)
		if token != "" && err == nil {
			return token, nil
		}
	}

	// 2. Demande du code d’appareil
	deviceCodeResp, err := obtainDeviceCode()
	if err != nil {
		return "", err
	}

	fmt.Println(deviceCodeResp.Message)

	// 3. Polling pour le token d’accès
	for {
		time.Sleep(time.Duration(deviceCodeResp.Interval) * time.Second)

		token, err := obtainAccessToken(deviceCodeResp.DeviceCode, false)
		if token != "" {
			return token, nil
		}
		if err != nil {
			return "", err
		}

		fmt.Println("En attente de l’authentification de l’utilisateur... Veuillez continuer sur votre navigateur.")
	}
}

func FindOutlookEvent(accessToken, subject, startDate, endDate string) (bool, error) {
	eventUrl := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/me/calendarView?startDateTime=%s&endDateTime=%s&$filter=subject%%20eq%%20'%s'",
		startDate, endDate, url.QueryEscape(subject),
	)

	req, err := http.NewRequest("GET", eventUrl, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("bad API HTTP status: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}

	if err = json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	events, exists := result["value"].([]interface{})
	if !exists {
		return false, nil
	}

	return len(events) > 0, nil
}

func CreateOutlookEvent(accessToken, subject, startDate, endDate string, isAllDay bool) error {
	eventUrl := "https://graph.microsoft.com/v1.0/me/events"
	event := map[string]interface{}{
		"subject": subject,
		"start": map[string]string{
			"dateTime": startDate,
			"timeZone": "Romance Standard Time",
		},
		"end": map[string]string{
			"dateTime": endDate,
			"timeZone": "Romance Standard Time",
		},
		"showAs":       "oof",
		"isReminderOn": false,
		"isAllDay":     isAllDay,
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

	if _, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("bad API HTTP status: " + resp.Status)
	}

	return nil
}

func obtainDeviceCode() (*DeviceCodeResp, error) {
	deviceCodeURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode?mkt=fr-FR", tenantID)
	resp, err := http.PostForm(deviceCodeURL, map[string][]string{
		"client_id": {clientID},
		"scope":     {"offline_access User.Read Calendars.ReadWrite"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceCodeResp DeviceCodeResp
	if err := json.Unmarshal(body, &deviceCodeResp); err != nil {
		return nil, err
	}
	return &deviceCodeResp, nil
}

func obtainAccessToken(code string, isRefreshToken bool) (string, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := url.Values{}
	data.Set("client_id", clientID)

	if isRefreshToken {
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", code)
	} else {
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		data.Set("device_code", code)
	}

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

	if resp.StatusCode != http.StatusOK {
		// Sur un mauvais statut, on ne renvoie pas d’erreur afin de laisser le polling se faire
		return "", nil
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	viper.Set("ms.access_token", tokenResp.AccessToken)
	viper.Set("ms.refresh_token", tokenResp.RefreshToken)
	err = viper.WriteConfig()
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
