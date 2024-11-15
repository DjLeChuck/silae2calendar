package silae

import (
	"errors"
	"strings"
	"time"

	"silae2calendar/ms"
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

type UserData struct {
	Firstname           string `json:"firstname"`
	Lastname            string `json:"lastname"`
	Trigram             string
	Token               string              `json:"token"`
	RefreshToken        string              `json:"refresh_token"`
	CurrentCollaborator CurrentCollaborator `json:"current_collaborator_base"`
}

func (ud *UserData) SetTrigram() {
	trigram := strings.ToUpper(string(ud.Firstname[0]))
	lastNameParts := strings.Fields(ud.Lastname)

	if len(lastNameParts) == 1 {
		trigram += strings.ToUpper(lastNameParts[0][:2])
	} else {
		trigram += strings.ToUpper(string(lastNameParts[0][0])) + strings.ToUpper(string(lastNameParts[1][0]))
	}

	ud.Trigram = trigram
}

type Request struct {
	Status string `json:"status"`
}

// Freedays DateStartDayType and DateEndDayType could have value: morning, afternoon or full
type Freedays struct {
	Abbr             string  `json:"abbreviation"`
	DateStart        string  `json:"date_start"`
	DateStartDayType string  `json:"date_start_day_type"`
	DateEnd          string  `json:"date_end"`
	DateEndDayType   string  `json:"date_end_day_type"`
	Request          Request `json:"request"`
}

func (f Freedays) DateStartForOutlook() (string, error) {
	d, err := time.Parse(time.RFC3339, f.DateStart)
	if err != nil {
		return "", err
	}

	switch f.DateStartDayType {
	case "morning":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			9, 0, 0, 0,
			d.Location(),
		)

		return nd.Format(ms.EventDateLayout), nil
	case "afternoon":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			14, 0, 0, 0,
			d.Location(),
		)

		return nd.Format(ms.EventDateLayout), nil
	case "full":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			0, 0, 0, 0,
			d.Location(),
		)

		return nd.Format(ms.EventDateLayout), nil
	}

	return "", errors.New("invalid date start day type")
}

func (f Freedays) DateEndForOutlook() (string, error) {
	d, err := time.Parse(time.RFC3339, f.DateEnd)
	if err != nil {
		return "", err
	}

	switch f.DateEndDayType {
	case "morning":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			12, 30, 0, 0,
			d.Location(),
		)

		return nd.Format(ms.EventDateLayout), nil
	case "afternoon":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			18, 0, 0, 0,
			d.Location(),
		)

		return nd.Format(ms.EventDateLayout), nil
	case "full":
		nd := time.Date(
			d.Year(), d.Month(), d.Day(),
			0, 0, 0, 0,
			d.Location(),
		)

		return nd.Add(24 * time.Hour).Format(ms.EventDateLayout), nil
	}

	return "", errors.New("invalid date end day type")
}

func (f Freedays) IsAllDay() bool {
	return f.DateStartDayType == "full"
}

type CollaboratorFreedays struct {
	Freedays []Freedays `json:"freedays"`
}

type FreedaysData struct {
	CollaboratorFreedays []CollaboratorFreedays `json:"collaborator_freedays"`
}

type Criteria interface{}

type DateRangeCriteria struct {
	Type string `json:"_type"`
	Min  string `json:"min"`
	Max  string `json:"max"`
}

type StringValueCriteria struct {
	Type  string `json:"_type"`
	Value int    `json:"value"`
}

type ListStringValueCriteria struct {
	Type   string `json:"_type"`
	Values []int  `json:"values"`
}

type Filter struct {
	Type     string   `json:"_type"`
	Name     string   `json:"name"`
	Criteria Criteria `json:"criteria"`
}

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type RequestPayload struct {
	Filters []Filter `json:"filters"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Sort    Sort     `json:"sort"`
}
